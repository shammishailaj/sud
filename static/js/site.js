$(function () {  
  window.sud = new Object();
  window.sud.Send = function(Url,Data,SuccessF,ErrorF){
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
  window.sud.Login = function(login, password, configurationName){
    var dataResult
    window.sud.Send("/json/login",{
        login: login,
        password: password,
        configurationName: configurationName,
      },function(data){        
        dataResult = data        
        console.log( "json login:", data);
      },function(data){
        console.log( "json login errpr:", data);
      })
  };
  window.sud.BeginTransaction = function(){
    var transactionUID 
    var error
    window.sud.Send("/json/begintransaction",{      
    },function(data){
      transactionUID = data.TransactionUID       
      error = data.error
    },function(data){
      error = "connection error"
    })
    return {transactionUID: transactionUID, error: error}
  };
  window.sud.CommitTransaction = function( transactionUID){
    var commit
    var error
    window.sud.Send("/json/committransaction",{      
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
  window.sud.RollbackTransaction = function( transactionUID){
    var rollback
    var error
    window.sud.Send("/json/rollbacktransaction",{      
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
  window.sud.Login("Test", "Test", "Storage.Default")
});
