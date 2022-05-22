


var page = require('webpage').create(),
// pretend to be a different browser, helps with some shitty browser-detection scripts
 // page.settings.userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36';
// page.settings.userAgent = 'Mozilla/5.0 (Windows NT 10.0; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0';

system = require('system');

page.viewportSize = { width: 1920, height: 1080 };
page.open(system.args[1],function (){
	console.log(page.content);
	phantom.exit();
});
