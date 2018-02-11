// Client side code
// This is quite an abomination of manual DOM manipulation


// Function wrapper
(function(){
    // Initialise by grabbing the DOM elements
    var form = document.getElementById('form-shorten-url');
    var fInputShort = form.querySelector('[name=short]');
    var fInputLong = form.querySelector('[name=long]');
    var fInputOverwrite = form.querySelector('[name=overwrite]');
    var dataTable = document.getElementById('current-urls-data');
    var dialog = document.getElementById('dialog');

    // Generate some DOM elements
    var noShortUrlsRow = document.createElement('tr');
    noShortUrlsRow.innerHTML = '<td colspan="3">No URLS found. Add some today!</td>';
    noShortUrlsRow.style.display = 'none';
    dataTable.appendChild(noShortUrlsRow);

    // Hides the dialog (currently on timer)
    function closeDialog() {
       dialog.style.display = 'none'; 
    }

    // Populate the table on page load
    // Quick row generating function which appends to table
    function generateTableRow(record) {
        var row = document.createElement('tr');
        var short = document.createElement('td');
        var long = document.createElement('td');
        var actions = document.createElement('td');
        var shortA = document.createElement('a');
        shortA.className = 'link';
        shortA.href = '/' + record.short;
        shortA.text = record.short;
        short.appendChild(shortA);
        var longA = document.createElement('a');
        longA.className = 'link';
        longA.text = record.long;
        longA.href = record.long;
        long.appendChild(longA);
        
        // Work on action buttons
        /*var copyButton = document.createElement('button'); 
        copyButton.className = 'pure-button button-copy'; 
        copyButton.innerHTML = 'cp';
        copyButton.addEventListener('click', function() {
            copyRecord(row);
        });*/
        var deleteButton = document.createElement('button'); 
        deleteButton.className = 'pure-button button-delete'; 
        deleteButton.innerHTML = 'Delete';
        deleteButton.addEventListener('click', function() {
            deleteRecord(row);
        });
        // actions.appendChild(copyButton);
        actions.appendChild(deleteButton);

        row.appendChild(short);
        row.appendChild(long);
        row.appendChild(actions);
        dataTable.appendChild(row);
    }

    function deleteRecord(ref) {
        var short = ref.children[0].children[0].text;
        fetch('/api/delete/' + short, {
            credentials: 'include',
            method: 'DELETE'
        }).then(function() {
            dataTable.removeChild(ref);
            if (dataTable.children.length <= 1) {
                noShortUrlsRow.style.display = '';
            }
            dialog.className = 'dialog dialog-success';
            dialog.innerHTML = 'Short URL <b>' + short + '</b> successfully created.';
            dialog.style.display = '';
        }).catch(function(e) {
            dialog.className = 'dialog dialog-error';
            dialog.innerHTML = 'An unknown error occurred, check the console for more details.';
            dialog.style.display = '';
            console.error(e);
        });
    }

    function copyRecord(ref) {
        var link = ref.children[0];
        link.select();
        document.execCommand('copy');
    }

    // Request using fetch
    fetch('/api/list', {credentials: 'include'}).then(function(data){return data.json()})
        .then(function(res) {
            // Use data to populate table
            dataTable.removeChild(dataTable.children[0]);
            // Needs to be synchronous
            for (var i = 0; i < res.length; i++) {
                generateTableRow(res[i]);
            }
            if (dataTable.children.length <= 1) {
                noShortUrlsRow.style.display = '';
            }
        })
        .catch(function(err) {
            alert("An unknown error occurred, check the console for more details");
            console.error(err);
        });

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
        .then(function(res){
            // Generate table row
            if (res.status) {
                if (data.overwrite) {
                    // Check to see if anything can be overwritten.
                    var found = false;
                    for (var i of dataTable.children) {
                        if (i.children.length == 1) {
                            continue;
                        }
                        if (i.children[0].children[0].text == res.short) {
                            found = true;
                            i.children[1].children[0].href = data.long;
                            i.children[1].children[0].text = data.long;
                            break;
                        }
                    }
                    if (!found) {
                        generateTableRow({
                            short: res.short,
                            long: data.long
                        });
                    }
                } else {
                    generateTableRow({
                        short: res.short,
                        long: data.long
                    });
                }
                dialog.className = 'dialog dialog-success';
                dialog.innerHTML = 'Short URL <b>' + res.short + '</b> successfully created.';
                dialog.style.display = '';
                dataTable.children[0].style.display = 'none';
                // Clear inputs
                fInputShort.value = '';
                fInputLong.value = '';
            } else {
                switch(res.error) {
                    case 'ExistsShortURL':   
                        dialog.innerHTML = 'Short URL already exists. Use overwrite to reuse the current name.';
                        break
                    case 'InvalidShortURL':
                        dialog.innerHTML = 'Short URL is invalid.'
                        break;
                    default:
                        dialog.innerHTML = 'An unknown error occurred, check the console for more details';
                        console.error(res.error);
                }   
                dialog.className = 'dialog dialog-error';
                dialog.style.display = '';
            }
                
        }).catch(function(err) {
            dialog.className = 'dialog dialog-error';
            dialog.innerHTML = 'An unknown error occurred, check the console for more details.';
            dialog.style.display = '';
            console.error(err);
        });
    });
})();
