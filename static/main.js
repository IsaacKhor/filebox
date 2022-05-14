window.onload = () => main()

let selectedFiles = new Set()
let tableRows = undefined

function main() {
    // Table selections
    tableRows = document.querySelectorAll('#files tbody tr')
    tableRows.forEach( (val, i, _) =>
            val.addEventListener('click', e => onTableRowClick(val, e)))

    // Uncheck checkboxes that remain checked through page refreshes
    document.querySelectorAll('#files tbody input.file-selector').forEach(
        (val, i, _) => {
            val.checked = false
        })

    document.getElementById('btn-download').addEventListener('click', downloadFiles)
    document.getElementById('btn-delete').addEventListener('click', deleteFiles)
    document.getElementById('btn-create-view').addEventListener('click', createView)
}

function onTableRowClick(row, event) {
    let fileid = row.dataset.fileid
    if(selectedFiles.has(fileid)) {
        selectedFiles.delete(fileid)
    } else {
        selectedFiles.add(fileid)
    }

    syncCheckboxStates()
}

function syncCheckboxStates() {
    tableRows.forEach((row, i, _) =>
        getCheckboxFromRow(row).checked = selectedFiles.has(row.dataset.fileid))
}

function getCheckboxFromRow(row) {
    return row.querySelector('input[type=checkbox]')
}

function downloadFiles() {
    alert("Not implemented")
}

function deleteFiles() {
    alert("Not implemented")
}

function createView() {
    alert("Not implemented")
}