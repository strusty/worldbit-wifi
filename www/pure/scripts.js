var host = '//95.216.200.100:3000';
var buttonPhoneTitle = 'Send code by SMS';
var buttonCodeTitle = 'Authorize';

document.cookie = '';

function getCookie(name) {
  var matches = document.cookie.match(new RegExp(
    "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
  ));
  return matches ? decodeURIComponent(matches[1]) : undefined;
}

function setCookie(name, value, options) {
  options = options || {};

  var expires = options.expires;

  if (typeof expires == "number" && expires) {
    var d = new Date();
    d.setTime(d.getTime() + expires * 1000);
    expires = options.expires = d;
  }
  if (expires && expires.toUTCString) {
    options.expires = expires.toUTCString();
  }

  value = encodeURIComponent(value);

  var updatedCookie = name + "=" + value;

  for (var propName in options) {
    updatedCookie += "; " + propName;
    var propValue = options[propName];
    if (propValue !== true) {
      updatedCookie += "=" + propValue;
    }
  }

  document.cookie = updatedCookie;
}

var stage = 0;

var tariff = [];

// stage: 0 - send code, 1 - registration, 2 - payment, 3 - successful fare selection

function changeStage(value) {
  allTariff(value);
  renderBackButton(value);
  renderInput(value);
  renderButton(value);
  renderForm(value);
  showPaypal(value);
  headerPayment(value);
  stage = value;
  setCookie('s', value)
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

  if (_stage === 0 && phone && code) {
    input.name = 'phone';
    input.placeholder = "Phone number +XXXXX";
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

function auth() {
  var code = form.elements['code'].value;
  if (code.length === 0) {
    setLoading(false);
    return setError('Please fill code');
  }
  methods.post(host + '/auth', "POST", { confirmationCode: code }, function (el) {
    changeStage(2)
  },);
}

function renderForm(_stage) {
  var loginForm = document.querySelector('.Login-form');
  if (_stage > 1 && form) {
    loginForm.removeChild(form);
  } else if (!form) {
    loginForm.appendChild(form);
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
        voucherError.innerHTML = 'Internet access provided!';
        window.location.replace('https:google.com')
      } else if (chilliController.clientState === 0) {
        voucherError.innerHTML = 'Error: Invalid voucher';
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
      setCookie('PhoneNumber', phoneNumber)
    }
    changeStage(1);
  });
}


function createBut(name, inner, currency, value) {
  const alertText = document.querySelector('.alert-block--text');
  const alertWrapper = document.querySelector(".alert-block--text-wrapper")


  const button = document.createElement('button');
  button.className = name;
  button.innerHTML = inner;
  button.onclick = function () {
    const alert = document.querySelector('.alert-show');
    alert.className = 'alert';
    const alertOk = document.querySelector('#alert-block--button-ok');
    const alertNo = document.querySelector('#alert-block--button-no');

    alertOk.onclick = function () {
      const PhoneNumber = getCookie('PhoneNumber');
      const alertWrapper = document.querySelector(".alert-block--text-wrapper")

      alertOk.style.display = "none";
      alertNo.style.display = "none";
      alertText.style.fontSize = '20px';
      alertText.innerHTML = 'Loading...';
      alertWrapper.style = 'flex: 1 1 0%;';
      methods.post(host + '/crypto/payment', "POST", {
          PhoneNumber: PhoneNumber,
          Currency: currency,
          PricingPlanID: value.id
        },
        function (el) {
          alertText.style.fontSize = '13px';
          if (el.address && el.amount) {
            const walletOk = document.querySelector('#alert-block--button-wallet');
            walletOk.className = 'alert-block--button-on';
            walletOk.style.display = '';
            alertWrapper.style = 'flex: 0 0 0%;';
            walletOk.onclick = function () {
              alert.className = 'alert-show';
              const alertInput = document.querySelector('.input-wrapper');
              alertInput.className = "input-wrapper-show";
              alertOk.style.display = "";
              alertNo.style.display = "";
              walletOk.className = 'alert-block--button-on-show';
            };
            const inputShow = document.querySelector('.input-wrapper-show');
            inputShow.className = 'input-wrapper';
            const inputAddress = document.querySelector('.input-alert-address');
            inputAddress.value = el.address;
            const inputSum = document.querySelector('.input-alert-sum');
            inputSum.value = el.amount;
            alertText.innerHTML = ''
          }
        })
    };
    alertWrapper.style = 'flex: 1 1 0%;';
    alertText.innerHTML = 'are you sure?';
    alertText.style.fontSize = 'xx-large';
    alertNo.onclick = function () {
      alert.className = 'alert-show';
    }
  };
  return button
}

function tariffSend() {
  if (getCookie('s') === '2') {
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
    var script = document.createElement('script');
    script.src = "https://www.paypalobjects.com/api/checkout.js";
    script.onload = function() {
      methods.post(host + '/plans', "GET", null, function (res) {
        tariff = res;
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
    };
    document.getElementsByTagName('head')[0].appendChild(script);
  }
}

