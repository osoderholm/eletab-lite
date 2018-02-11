var Eletab = function () {
    this.token = "";
};

Eletab.prototype.apiCall = function(path, data, success, error) {
    var settings = {
        "async": true,
        "crossDomain": true,
        "url": "/api/v1"+path,
        "method": "GET",
        "headers": {
            "authorization": "Bearer "+this.token,
            "cache-control": "no-cache"
        },
        "data": data
    };

    $.ajax(settings).done(function(data) {
        success(data);
    }).fail(function(data) {
        error(data);
    });
};

Eletab.prototype.parsePayload = function (payload) {
    try{
        this.token = JSON.parse(payload).token;
        return this.token !== "";
    } catch (err) {
        return false;
    }
};

Eletab.prototype.destroy = function () {
    this.token = "";
};

Eletab.prototype.getAccount = function (success, error) {
    this.apiCall("/account/get", {}, function (data) {
        success(data);
    }, function (data) {
        error(data);
    })
};

Eletab.prototype.getAccounts = function (success, error) {
    this.apiCall("/account/get_all", {}, function (data) {
        success(data);
    }, function (data) {
        error(data);
    })
};

Eletab.prototype.addAccount = function (name, username, password, balance, level, success, error) {
    var stuff = {
        "name": name,
        "username": username,
        "pass": password,
        "balance": balance,
        "level": level
    };
    this.apiCall("/account/new", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.increaseBalance = function (id, sum, success, error) {
    var stuff = {
        "id": id,
        "sum": sum
    };
    this.apiCall("/account/increase", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.decreaseBalance = function (id, sum, success, error) {
    var stuff = {
        "id": id,
        "sum": sum
    };
    this.apiCall("/account/decrease", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.getCards = function (username, success, error) {
    var stuff = {
        "username": username
    };
    this.apiCall("/account/get_cards", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.addCard = function (username, cardId, success, error) {
    var stuff = {
        "username": username,
        "card_id": cardId
    };
    this.apiCall("/account/add_card", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.removeCard = function (username, cardId, success, error) {
    var stuff = {
        "username": username,
        "card_id": cardId
    };
    this.apiCall("/account/delete_card", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.getClients = function (success, error) {
    this.apiCall("/clients/get_all", {}, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.addClient = function (description, level, success, error) {
    var stuff = {
        "description": description.trim(),
        "level": level
    };
    this.apiCall("/clients/new", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

Eletab.prototype.removeClient = function (key, success, error) {
    var stuff = {
        "key": key.trim()
    };
    this.apiCall("/clients/delete", stuff, function (data) {
        success(data);
    }, function (data) {
        error(data);
    });
};

