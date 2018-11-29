var host = '//95.216.200.100:3000';
var buttonPhoneTitle = 'Send code by SMS';
var buttonCodeTitle = 'Authorize';

var stage = parseFloat(localStorage.getItem('s')) || 0;
var tariff = [];

// stage: 0 - send code, 1 - registration, 2 - payment, 3 - successful fare selection

function changeStage(value) {
  allTariff(value);
  renderBackButton(value);
  renderInput(value);
  renderButton(value);
  renderForm(value);
  showCaptcha(value);
  showPaypal(value);
  headerPayment(value);
  stage = value;
  localStorage.setItem('s', value);
}

var backButton = document.createElement('span');
backButton.className = 'back';
backButton.onclick = onBackClick;
backButton.innerHTML = 'change phone number';

function renderBackButton(_stage) {
  var back = document.querySelector(".back");
  if (!back && _stage === 1) {
    extra.appendChild(backButton);
  } else if (back && _stage !== 1) {
    extra.removeChild(backButton);
  }
}

function renderInput(_stage) {
  var phone = document.querySelector('input[name=phone]');
  var code = document.querySelector('input[name=code]');
  input.value = '';

  if (_stage === 0 && !phone && code) {
    input.name = 'phone';
    input.type = "number";
    input.placeholder = "Phone number...";
  } else if (_stage === 1 && phone && !code) {
    input.name = 'code';
    input.type = "default";
    input.placeholder = "Code...";
  }
}

function renderButton(_stage) {
  if (_stage === 0 && submitButton.id === 'code-btn') {
    submitButton.innerHTML = buttonPhoneTitle;
    submitButton.id = 'phone-btn';
  } else if (_stage === 1 && submitButton.id === 'phone-btn') {
    submitButton.innerHTML = buttonCodeTitle;
    submitButton.id = 'code-btn';
  }
}

function renderForm(_stage) {
  var loginForm = document.querySelector('.Login-form');
  if (_stage > 1 && form) {
    loginForm.removeChild(form);
  } else if (!form) {
    loginForm.appendChild(form);
  }
}

function showCaptcha(_stage) {
  var recaptchaContainer = document.getElementById('captcha');
  if (recaptchaContainer) {
    var ind = recaptchaContainer.className.indexOf(' show');
    if (_stage === 1 && ind < 0) {
      recaptchaContainer.className = recaptchaContainer.className + ' show';
    } else if (_stage !== 1 && ind >= 0) {
      recaptchaContainer.className = recaptchaContainer.className.slice(0, ind) + recaptchaContainer.className.slice(ind + 5);
    }
  }
}


function showPaypal(_stage) {
  var paypalContainer = document.getElementsByClassName('switch')[0];
  var ind = paypalContainer.className.indexOf(' show');
  if (_stage === 2 && ind < 0) {
    paypalContainer.className = paypalContainer.className + ' show';
  } else if (_stage !== 2 && ind >= 0) {
    paypalContainer.className = paypalContainer.className.slice(0, ind) + paypalContainer.className.slice(ind + 5);
  }
}

// then define event handler functions

var voucherError;

function sendVoucher() {
  const input = document.querySelector('input.voucher--input');
  const button = document.querySelector('button.voucher-button');
  voucherError.className = 'loading';
  voucherError.innerHTML = 'loading...';
  chilliController.logon(input.value, input.value);
}

function headerPayment(value) {
  if (value === 2) {
    document.querySelector('.block-switch').style.display = 'flex'
  }
}

var extra;
var form;
var submitButton;
var input;
document.onreadystatechange = function () {
  if (document.readyState == "interactive") {
    extra = document.querySelector(".extra");
    form = document.querySelector('form');
    submitButton = document.querySelector('button[type=submit]');
    submitButton.onclick = onSubmit;
    input = document.querySelector('#main-input');
    const buttonVoucher = document.querySelector('button.voucher-button');
    buttonVoucher.onclick = sendVoucher;
    input.addEventListener('input', function () {
      if (input.name === 'phone') {
        input.value = input.value.replace(/[^\d]/g, '');
        if (input.value.indexOf('+') !== 0) {
          input.value = "+" + input.value;
          if (input.value.length === 1) {
            input.value = ''
          }
        }
        setError();
      }
    });

    // connect chilli library
    voucherError = document.querySelector('div.voucher--error');

    function updateUI(cmd) {
      if (chilliController.clientState === 1) {
        window.location.replace('https:google.com')
      }
      voucherError.innerHTML = 'You called the method' + cmd +
        '\n Your current state is =' + chilliController.clientState;
    }

    function handleErrors(code) {
      voucherError.className = 'voucher--error';
      voucherError.innerHTML = 'Error: ' + code;
    }

    chilliController.host = "192.168.182.1";
    chilliController.port = "3990";
    chilliController.onError = handleErrors;
    chilliController.onUpdate = updateUI;
    changeStage(stage);
  }
};

var loader = document.createElement('span');
loader.className = 'loader';
loader.innerHTML = 'Loading...';

function setLoading(value) {
  var loaderEl = document.querySelector(".loader");
  if (value && !loaderEl) {
    extra.appendChild(loader);
  } else if (!value && loaderEl) {
    extra.removeChild(loader);
  }
}

var error = document.createElement('span');
error.className = 'error';

function setError(err) {
  var errorEl = document.querySelector(".error");
  if (err && !errorEl) {
    error.innerHTML = err;
    extra.appendChild(error);
  } else if (!err && errorEl) {
    extra.removeChild(error);
  }
}

function onSubmit() {
  setLoading(true);
  setError();
  if (stage === 0) {
    setTimeout(authSendCode, 1000);
  } else if (stage === 1) {
    setTimeout(auth, 1000);
  } else {
    setLoading(false);
    setError('Something went wrong! Please try again later.');
  }
  return false;
}

function onBackClick() {
  changeStage(0);
  return false;
}

/* api methods */

var methods = {
  post: function (_url, method, _body, cb) {
    var body = JSON.stringify(_body);
    var xhr = new XMLHttpRequest();
    xhr.open(method, _url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
      if (xhr.readyState === XMLHttpRequest.DONE) {
        if (xhr.status === 200) {
          var response = JSON.parse(xhr.response);
          if (response.error) {
            setLoading(false);
            setError(response.error.message);
          } else {
            setLoading(false);
            if (cb) cb(response);
          }
        } else {
          setLoading(false);
          setError('Something went wrong!');
        }
      }
    };
    xhr.send(body);
  }
};

function authSendCode() {
  var phoneNumber = form.elements['phone'].value;
  if (!phoneNumber) {
    setLoading(false);
    return setError('Please fill phone number');
  }
  methods.post(host + '/auth/sendCode', "POST", { phoneNumber: phoneNumber }, function (el) {
    if (el.success === true) {
      localStorage.setItem('PhoneNumber', phoneNumber)
    }
    changeStage(1);
  });
}

// data-verifycallbackname
function auth() {
  var captchaRes = grecaptcha.getResponse();
  var code = form.elements['code'].value;
  if (code.length === 0) {
    setLoading(false);
    return setError('Please fill code');
  }
  if (captchaRes.length === 0) {
    setLoading(false);
    return setError('Please check reCAPTCHA');
  }
  methods.post(host + '/auth', "POST", { confirmationCode: code, captcha: captchaRes }, function (el) {
    console.log('capch', el);
    changeStage(2)
  },);
}

// data-verifycallbackname
function cryptoPayment() {
  var captchaRes = grecaptcha.getResponse();
  var code = form.elements['code'].value;
  if (code.length === 0) {
    setLoading(false);
    return setError('Please fill code');
  }
  if (captchaRes.length === 0) {
    setLoading(false);
    return setError('Please check reCAPTCHA');
  }
  setLoading(false);
  changeStage(2);
}

function createBut(name, inner, currency, value) {
  const button = document.createElement('button');
  button.className = name;
  button.innerHTML = inner;
  button.onclick = function () {
    const alert = document.querySelector('.alert-show');
    alert.className = 'alert';
    const alertOk = document.querySelector('#alert-block--button-ok');
    alertOk.onclick = function () {
      alert.className = 'alert-show';
      const PhoneNumber = localStorage.getItem('PhoneNumber');
      methods.post(host + '/crypto/payment', "POST", {
          PhoneNumber: PhoneNumber,
          Currency: currency,
          PricingPlanID: value.id
        },
        function (el) {
          if (el.address && el.amount) {
            const address = document.querySelector('.wallet');
            address.innerHTML = `address - ${el.address}<br/> sum - ${el.amount}`
          }
        })
    };
    const alertNo = document.querySelector('#alert-block--button-no');
    alertNo.onclick = function () {
      alert.className = 'alert-show';
    }
  };
  return button
}

function tariffSend() {
  if (localStorage.getItem('s') === '2') {

    allTariff()
  }
}

tariffSend();

function div(name, text, teg, id) {
  const el = document.createElement(teg);
  el.className = name;
  el.innerHTML = text;
  el.id = id;
  return el
}

function allTariff(_state) {

  if (_state === 2) {

    const tariffChange = methods.post(host + '/plans', "GET", null, function (res) {
      tariff = res;
      console.log(3, tariff);
      res.map(function (value, index) {
        const newLi = document.createElement('li');
        newLi.className = 'list';
        const newDiv = document.createElement('div');
        newDiv.className = 'c-cont';
        const block = document.getElementById("crypto");
        newDiv.appendChild(div('c-h', `tariff-${index + 1}`, 'div'));
        newDiv.appendChild(div('c-price', "duration - " + value.amountUSD + " USD", 'div'));
        newDiv.appendChild(div('c-lim', "maxUsers - " + value.maxUsers, 'div'));
        newDiv.appendChild(div('c-lim', "upLimit - " + value.upLimit + " Gb", 'div'));
        newDiv.appendChild(div('c-lim', "downLimit - " + value.downLimit + " Gb", 'div'));
        newDiv.appendChild(div('c-lim', "purgeDays - " + value.purgeDays + " Days", 'div'));
        const divButton = div('butBlock', null, 'div');
        newDiv.appendChild(divButton);
        const buttonETH = createBut('button-buy', 'Buy for ETH', 'ETH', value);
        const buttonBTC = createBut('button-buy', 'Buy for BTC', 'BTC', value);
        divButton.appendChild(buttonETH);
        divButton.appendChild(buttonBTC);
        const divPaypal = div('paypal-button-', null, 'div', `paypal-button-${value.id}`);
        newDiv.appendChild(divPaypal);
        newLi.appendChild(newDiv);
        block.appendChild(newLi);
        pay(`paypal-button-${value.id}`, String(value.amountUSD), value.id)
      });
    });
  }
}

