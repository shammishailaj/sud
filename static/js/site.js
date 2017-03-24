$(function () {  
  window.ClientSud = function(){
    var Send = function(Url,Data,SuccessF,ErrorF){
      $.ajax({
        dataType: "json",
        url: Url,
        type: "POST",
        contentType: 'application/json; charset=utf-8',
        dataType: 'json',
        async: false,
        data: JSON.stringify(Data),
        success: SuccessF,
        error : ErrorF,
      });
    }
    var Client = new Object;
    Client.Login = function(login, password, configurationName){
      var error
      var loginOk
      Send("/json/login",{
          login: login,
          password: password,
          configurationName: configurationName,
      },function(data){
        loginOk = data.login       
        error = data.error
      },function(data){
        error = "connection error"
      })
      return {login: loginOk, error: error}
    };
    Client.BeginTransaction = function(){
      var transactionUID 
      var error
      Send("/json/begintransaction",{      
      },function(data){
        transactionUID = data.TransactionUID       
        error = data.error
      },function(data){
        error = "connection error"
      })
      return {transactionUID: transactionUID, error: error}
    };
    Client.CommitTransaction = function( transactionUID){
      var commit
      var error
      Send("/json/committransaction",{      
        transactionUID: transactionUID
      },function(data){
        commit = data.commit       
        error = data.error
      },function(data){
        commit = false
        error = "connection error"
      })
      return {commit: commit, error: error}
    };
    Client.RollbackTransaction = function(transactionUID){
      var rollback
      var error
      Send("/json/rollbacktransaction",{      
        transactionUID: transactionUID
      },function(data){
        rollback = data.rollback       
        error = data.error
      },function(data){
        rollback = false
        error = "connection error"
      })
      return {rollback: rollback, error: error}
    };  
    return Client;
  }  
  window.sud.Login("Test", "Test", "Storage.Default")
});
