function removeFromDb(cname, disease_code){
    fetch(`/record?cname=${cname}&disease_code=${disease_code}`, {method: "DELETE"}).then(res =>{
        if (res.status == 200){
            window.location.pathname = "/record"
        }
    })
 }
 
 function updateDb(cname, disease_code, total_patients, total_deaths) {
    let input1 = document.getElementById(cname)
    let input2 = document.getElementById(disease_code)
    let input3 = document.getElementById(total_patients)
    let input4 = document.getElementById(total_deaths)

    let new_cname = input1.value
    let new_disease_code = input2.value
    let new_total_patients = input3.value
    let new_total_deaths = input4.value

    fetch(`/record?cname=${new_cname}&disease_code=${new_disease_code}&total_patients=${new_total_patients}&total_deaths=${new_total_deaths}`, {method: "PUT"}).then(res =>{
        if (res.status == 200){
        alert("Database updated")
            window.location.pathname = "/record"
        }
    })
 }

 function createTable() {
    let input = document.getElementById(Records)
    each(input, function(i, f) {
        var tblRow = "<tr>" + "<td>" + f.cname + "</td>" +
        "<td>" + f.lastName + "</td>" + "<td>" + f.job + "</td>" + "<td>" + f.roll + "</td>" + "</tr>"
        $(tblRow).appendTo("#record tbody");
    });

}