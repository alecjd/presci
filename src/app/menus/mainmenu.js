const {Menu} = require('electron')
const electron = require('electron')
const dialog = require('electron').dialog;
const app = electron.app

function openPrefetchFile () {
    const prefetchFiles = dialog.showOpenDialog({
        properties: [
            'openFile'
        ],
        filters: [
            {
                name: 'prefetch',
                extensions: [
                    'pf'
                ]
            }
        ]
    })
    if (!prefetchFiles) return
    console.log(prefetchFiles[0])
};

const template = [
    {
        'label': 'File',
        'submenu': [
            {
                'label': 'Open File',
                click () {
                    openPrefetchFile();
                }
            }
        ]
    }
]

const menu = Menu.buildFromTemplate(template)
Menu.setApplicationMenu(menu)