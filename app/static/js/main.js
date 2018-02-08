$(document).ready(function () {
    var loginManager = new AccountLoginManager();

    var user_username   = $("#user_username");
    var user_password   = $("#user_password");

    $("#user_btn_login").click(function () {
        var username = user_username.val();
        var password = user_password.val();

        loginManager.login(username, password,
            function (data) {
                console.log("Success!");
                var token = data.token;
                console.log(token);
                var payload = btoa(JSON.stringify(data));
                localStorage.setItem("payload", payload);
                $(location).attr("href", "/user.html");
            }, function () {
                alert("Login failed!");
                user_username.val("");
                user_password.val("");
            });


    });

    var admin_username  = $("#admin_username");
    var admin_password  = $("#admin_password");

    $("#admin_btn_login").click(function () {
        var username = admin_username.val();
        var password = admin_password.val();

        loginManager.login(username, password,
            function (data) {
                console.log("Success!");
                var token = data.token;
                console.log(token);
                var payload = btoa(JSON.stringify(data));
                localStorage.setItem("payload", payload);
                $(location).attr("href", "/admin.html");
            }, function () {
                alert("Login failed!");
                admin_username.val("");
                admin_password.val("");
            });
    });

});