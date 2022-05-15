


var page = require('webpage').create(),
// pretend to be a different browser, helps with some shitty browser-detection scripts
 // page.settings.userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36';
// page.settings.userAgent = 'Mozilla/5.0 (Windows NT 10.0; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0';

system = require('system');
var imageIndex = 0;

page.onResourceRequested = function(request) {
	if (request.url.indexOf('www.google.com/search?tbm=map&fp=')>-1)
		console.log(request.url);
};

page.onLoadFinished = function(status) {
	//console.log(page.url);//Debug
	setTimeout(processPage, 5000);
};

page.viewportSize = { width: 1920, height: 1080 };
page.open(system.args[1]);


var processPage = function(){
	// page.render('google_map_page-'+imageIndex+'.png');//Debug
	if (imageIndex==0)
	{
		console.log(page.content);
		console.log("<<FirstPageEnd>>");
	}
	
	//Go next page
	imageIndex++;
	var btn = page.evaluate(function() {
		var btn = document.getElementById('section-pagination-button-next');
		btn.click();
		return btn.outerHTML;
	});
	//console.log(btn);//Debug
	

	if (btn==null)			
		phantom.exit();
	
	var pos = btn.indexOf('disabled="true"');
	//console.log(pos);//Debug
	if (pos>-1 || imageIndex>5)
		phantom.exit();		
	else
		setTimeout(processPage, 5000);		
};



