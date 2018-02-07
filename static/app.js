// Function wrapper
(function(){
    // Initialise by grabbing the DOM elements
    var form = document.getElementById('form-shorten-url');
    var fInputShort = form.querySelector('[name=short]');
    var fInputLong = form.querySelector('[name=long]');
    var fInputOverwrite = form.querySelector('[name=overwrite]');
    var dataTable = document.getElementById('current-urls-data');

    // Generate some DOM elements
    var noShortUrlsRow = document.createElement('tr');
    noShortUrlsRow.innerHTML = '<td colspan="3">No URLS found. Add some!</td>'

    // Populate the table on page load
    

    // Attach an event listener to the form 
    // (which submits the form and appends it to table)
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        var data = {
            short: fInputShort.value,
            long: fInputLong.value,
            overwrite: fInputOverwrite.checked
        };
        fetch('/api/shorten', {
            credentials: 'include',
            method: 'post',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data) 
        }).then(function(res){return res.json()})
        .then(function(data){
            console.log(data);
        });
    });
})();
