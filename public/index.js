function removeFromDb(item){
    fetch(`/record?item=${item}`, {method: "Delete"}).then(res =>{
        if (res.status == 200){
            window.location.pathname = "/record"
        }
    })
 }
 
 function updateDb(item) {
    let input = document.getElementById(item)
    let newitem = input.value
    fetch(`/record?olditem=${item}&newitem=${newitem}`, {method: "PUT"}).then(res =>{
        if (res.status == 200){
        alert("Database updated")
            window.location.pathname = "/record"
        }
    })
 }