#Wi-Fi crystalline front-end

This is Captive portal based on CoovaChilli. Main purpose - Wi-Fi distribution.
For receiving Internet via Wi-Fi you need to connect to Crystalline Wi-Fi hotspot,
log in and pay acccording to suitable tariff plan. Payment can be made two ways:
either PayPal service, or Cryptocurrency.
This device is using OpenWrt 18.06.1 r7258-5eb055306f firmware. It can work directly with Captive Portal,
made with [CoovaChili](https://coova.github.io/) (which is the heir of [Chillispot](http://www.chillispot.org/))
portal is connected with FreeRADIUS server which sets directions for CoovaChili.
FreeRADIUS is a server with open-source code, responsible for authorisation and autentification;
Frontend developed with JavaScript, HTML, CSS;
Also ChilliLibrary was used for froen-end communication with CoovaChili and PayPal for tariff plans payments.



##Start in dev environment
####Project is using YARN package manager.

Download dependencies via YARN 

```yarn```

Start project

```yarn start```

###Start in production

 1) Install firmware OPENWrt to devices;
 2) Install CoovaChili Captive Portal;
 3) Deploy FreeRADUIS;
 4) Configure CoovaChili config so he can communicate with FreeRADIUS и OPENWrt
 5) In CoovaChilli setting need to specify IP address that contains Wi-Fi crystalline frontend

