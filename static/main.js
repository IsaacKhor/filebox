window.onload = () => main()

let selectedFiles = new Set()
let tableRows = undefined
let opsInProgress = 0
let uploadXhr = undefined

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
    document.getElementById('btn-upload-file').addEventListener('click', uploadFiles)
    document.getElementById('btn-cancel-upload').addEventListener('click', cancelUpload)
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

function cancelUpload(ev) {
    ev.preventDefault()
    uploadXhr.abort()
    uploadXhr = null
}

function uploadFiles(ev) {
    ev.preventDefault()

    // Cancel any previous uploads
    if (uploadXhr) uploadXhr.abort()

    let progressBar = document.querySelector('progress#file-upload-progbar')
    let progText = document.querySelector('span#file-upload-progtext')
    let xhr = newXhrForPath("POST", "/files")
    let fd = new FormData(document.querySelector('form#file-upload-form'))

    console.info(`Uploading ${fd.getAll('files').length} files.`)

    xhr.upload.addEventListener('progress', ev => {
        progressBar.value = ev.loaded / ev.total * 100
        progText.textContent = `Progress: ${ev.loaded}/${ev.total}`
    })
    xhr.upload.addEventListener('load', ev => window.location.reload())
    xhr.upload.addEventListener('error', ev => {
        console.error(ev)
        alert("An error occurred during file upload. Check logs for details")
    })
    uploadXhr = xhr
    xhr.send(fd)
}

function downloadFiles() {
    alert("Not implemented")
}

function deleteFiles() {
    console.log(`Deleting ${selectedFiles}`)
    for(let f of selectedFiles) {
        let xhr = newXhrForPath("DELETE", "files/" + f)
        opsInProgress++
        xhr.onload = event => {
            opsInProgress--
            if(opsInProgress === 0) {
                window.location.reload()
            }
        }
        console.log(`XHR DELETE request for ${f}`)
        xhr.send()
    }
}

function newXhrForPath(method, path) {
    let req = new XMLHttpRequest()
    req.open(method, path)
    return req
}

function createView() {
    alert("Not implemented")
}