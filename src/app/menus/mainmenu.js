const {Menu} = require('electron')
const electron = require('electron')
const dialog = require('electron').dialog;
const app = electron.app

var showOpen = function () {
    dialog.showOpenDialog({
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
};

const template = [
    {
        'label': 'File',
        'submenu': [
            {
                'label': 'Open File',
                click () {
                    showOpen();
                }
            }
        ]
    }
]

const menu = Menu.buildFromTemplate(template)
Menu.setApplicationMenu(menu)