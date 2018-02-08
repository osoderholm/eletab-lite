var eletab = new Eletab();
var level = 0;

window.onbeforeunload = function() {
    //localStorage.removeItem("payload");
    return '';
};

function getAccount() {
    eletab.getAccount(function (data) {
        if (data.level == 0){
            goHome();
            return
        }
        level = parseInt(data.level);
        $("#account_name").html(data.name);
        $("#account_username").html(data.username);
        $("#account_level").html(data.level);
    }, function () {
        alert("Could not get account info.");
    })
}

function getAccounts() {
    var listTemplate = '<% data.forEach(function(account) { %> \
                        <div class="row">\
                            <span><%- account.id %></span>\
                            <span><%- account.name %></span>\
                            <span><%- account.username %></span>\
                            <span><%- account.balance %></span>\
                            <span><%- account.level %></span>\
                            <span><%- account.disabled %></span>\
                            <span><button onclick="increaseBalance(<%- account.id %>)" >+ €€€</button></span>\
                            <span><button onclick="decreaseBalance(<%- account.id %>)" >- €€€</button></span>\
                            <span><button onclick="addCard(\'<%- account.username %>\')" >Add card</button></span>\
                            <span><button onclick="removeCard(\'<%- account.username %>\')" >Remove card</button></span>\
                            <span><button onclick="editAccount(\'<%- account.username %>\')" >Edit</button></span>\
                            <span><button onclick="deleteAccount(\'<%- account.username %>\')" >Delete</button></span>\
                        </div>\
                    <% }) %> \
                    ';
    var template = _.template(listTemplate);
    eletab.getAccounts(function (data) {
        $("#accounts_container").html(template({data: data}));
    }, function () {
        alert("Failed to get accounts");
    });
}

function increaseBalance(accountId) {
    var sumTxt = prompt("Enter amount in cents to ADD to balance (1€->100)", "");
    if (sumTxt != null && sumTxt != "") {
        var sum = parseInt(sumTxt);
        if(confirm("ADD " + sum + " to balance?")){
            eletab.increaseBalance(accountId, sum, function (data) {
                alert("Transaction was "+(data.accepted? "successful" : "not successful")+"!\n\n"+JSON.stringify(data));
                getAccounts();
            }, function () {
                alert("Insertion failed!");
            });
        }
    }
}

function decreaseBalance(accountId) {
    var sumTxt = prompt("Enter amount in cents to REMOVE from balance (1€->100)", "");
    if (sumTxt != null && sumTxt != "") {
        var sum = parseInt(sumTxt);
        if(confirm("REMOVE " + sum + " from balance?")){
            eletab.decreaseBalance(accountId, sum, function (data) {
                alert("Transaction was "+(data.accepted? "successful" : "not successful")+"!\n\n"+JSON.stringify(data));
                getAccounts();
            }, function () {
                alert("Remove failed!");
            });
        }
    }
}

function addCard(username) {
    var cardId = prompt("Add card to user '"+username+"'\nEnter card ID", "");
    if (cardId != null && cardId != "") {
        if (confirm("Add card with ID '"+cardId+"' to user '"+username+"'?")) {
            eletab.addCard(username, cardId, function (data) {
                alert("Card added!\nCard ID:"+data.card_id+"\nAccount ID:"+data.account.id);
            }, function () {
                alert("Adding card failed!");
            });
        }
    }
}

function removeCard(username) {
    eletab.getCards(username, function (data) {
        if (data == null) {
            alert("User '"+username+"' has no cards!");
            return
        }
        var cardsList = "";
        data.forEach(function (card) { cardsList += "\nCard ID: "+card.card_id });
        var removeID = prompt("Cards for user '"+username+"'"+cardsList+"\n\nEnter Card ID of card to remove:");
        if (removeID != null && removeID != "") {
            if(confirm("Are you sure?")) {
                eletab.removeCard(username, removeID, function (data) {
                    alert("Card with card ID '"+removeID+"' "+(data? "removed" : "NOT removed")+" from user '"+username+"'");
                }, function () {
                    alert("Remove card failed!");
                });
            }
        }

    }, function () {
        alert("Could not get cards for user '"+username+"'");
    });
}

function getClients() {
    var listTemplate = '<div class="row header">\
                            <span>ID</span>\
                            <span>Description</span>\
                            <span>Key</span>\
                            <span>Secret</span>\
                            <span>Delete</span>\
                        </div>\
                        <% data.forEach(function(client) { %> \
                        <div class="row">\
                            <span><%- client.id %></span>\
                            <span><%- client.description %></span>\
                            <span><%- client.api_key %></span>\
                            <span><%- client.secret %></span>\
                            <span><%- client.level %></span>\
                            <span><button onclick="deleteClient(\'<%- client.api_key %>\', \'<%- client.description %>\')" >Delete</button></span>\
                        </div>\
                    <% }) %> \
                    ';
    var template = _.template(listTemplate);
    eletab.getClients(function (data) {
        if(data != null) {
            $("#clients_container").html(template({data: data}));
        }
    }, function () {
        alert("Failed to get clients");
    });
}

function addClient() {
    var description = $("#new_client_description").val();
    var level = $("#new_client_level").val();
    eletab.addClient(description, level, function (data) {
        alert("Added client!\nDescription: "+
            data.description+"\nKey: "+
            data.api_key+"\nSecret: "+
            data.secret+"\nBalance: "+
            data.balance+"\nLevel: "+
            data.level);
        $("#new_client_description").val("");
        $("#new_client_level").val(1);
        getClients();
    }, function (data) {
        var str = "";
        data.forEach(function (value) { str = str+"\n"+value.error });
        alert("Errors:\n"+str);
    });
}

function deleteClient(key, description) {
    if(confirm("Are you sure you want to remove client with key '"+key+"' and description '"+description+"'")){
        eletab.removeClient(key, function (data) {
            alert("Client with description '"+description+"' "+(data? "deleted":"NOT deleted")+"!")
        }, function () {
            alert("Delete client failed!")
        });
    }

}

function setupListeners(){
    $("#btn_add_account").click(function () {
        var name = $("#new_name").val();
        var username = $("#new_username").val();
        var password = $("#new_password").val();
        var balance = $("#new_balance").val();
        var level = $("#new_level").val();
        eletab.addAccount(name, username, password, balance, level,
            function (data) {
                alert("Added account!\nName: "+
                    data.name+"\nUsername: "+
                    data.username+"\nBalance: "+
                    data.balance+"\nLevel: "+
                    data.level);
                $("#new_name").val("");
                $("#new_username").val("");
                $("#new_password").val("");
                $("#new_balance").val("");
                $("#new_level").val(0);
                getAccounts();
            }, function (data) {
                var str = "";
                data.forEach(function (value) { str = str+"\n"+value.error });
                alert("Errors:\n"+str);
            });
    });

    $("#btn_refresh_accounts").click(function () {
        getAccounts();
    });

    $("#btn_add_client").click(function () {
       addClient();
    });
}

$(document).ready(function () {
    var payload = localStorage.getItem("payload");
    if (!eletab.parsePayload(atob(payload))){
        $(location).attr("href", "/");
        return
    }

    setupListeners();

    getAccount();

    getAccounts();

    getClients();
});