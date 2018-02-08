var eletab = new Eletab();

window.onbeforeunload = function() {
    localStorage.removeItem("payload");
    return '';
};

function goHome() {
    $(location).attr("href", "/");
}

function getAccount() {
    eletab.getAccount(function (data) {
        $("#account_name").html(data.name);
        $("#account_username").html(data.username);
        var balance = parseFloat(data.balance)/100.00;
        $("#account_balance").html(balance);
    }, function () {
        alert("Could not get account info.");
    })
}

$(document).ready(function () {
    var payload = localStorage.getItem("payload");
    if (!eletab.parsePayload(atob(payload))){
        goHome();
        return
    }

    $("#account_refresh").click(function () {
        getAccount();
    });

    getAccount();
});