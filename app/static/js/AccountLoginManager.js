AccountLoginManager = function () {
    this.loginPath = "/api/v1/account_login";
};

AccountLoginManager.prototype.login = function (username, password, success, error) {
    var userData = JSON.stringify({
        "username": username.trim(),
        "password": password
    });
    var settings = {
        "async": true,
        "crossDomain": false,
        "url": this.loginPath,
        "method": "POST",
        "headers": {
            "cache-control": "no-cache"
        },
        "data": userData,
        "dataType": "json"
    };

    $.ajax(settings).done(function(data) {
        success(data);
    }).fail(function(data) {
        error(data);
    });
};