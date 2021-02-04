fetch('http://localhost:8888/graphql', 
    {
        method: 'POST', 
        headers: {'Content-type':'application/json', 'Origin': "http://fakewebsite.com"},
        body: JSON.stringify({operationName: "HaslettBankStatements", query: "{readBankStatements{Category\nPercentage\nTotal}}"})
    })
.then(response => response.json())
.then(data => console.log(data))