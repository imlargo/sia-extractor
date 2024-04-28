const chunk1 = require('./artifacts/chunk-1.json');
const chunk2 = require('./artifacts/chunk-2.json');
const chunk3 = require('./artifacts/chunk-3.json');
const chunk4 = require('./artifacts/chunk-4.json');

const fs = require('fs');

const keysMinas = [
    "3515 INGENIERÍA ADMINISTRATIVA", "3528 INGENIERÍA ADMINISTRATIVA", "3527 INGENIERÍA AMBIENTAL", "3529 INGENIERÍA AMBIENTAL", "3516 INGENIERÍA CIVIL", "3530 INGENIERÍA CIVIL", "3517 INGENIERÍA DE CONTROL", "3531 INGENIERÍA DE CONTROL", "3518 INGENIERÍA DE MINAS Y METALURGIA", "3532 INGENIERÍA DE MINAS Y METALURGIA", "3519 INGENIERÍA DE PETRÓLEOS", "3533 INGENIERÍA DE PETRÓLEOS", "3520 INGENIERÍA DE SISTEMAS E INFORMÁTICA", "3534 INGENIERÍA DE SISTEMAS E INFORMÁTICA", "3521 INGENIERÍA ELÉCTRICA", "3535 INGENIERÍA ELÉCTRICA", "3522 INGENIERÍA GEOLÓGICA", "3536 INGENIERÍA GEOLÓGICA", "3523 INGENIERÍA INDUSTRIAL", "3537 INGENIERÍA INDUSTRIAL", "3524 INGENIERÍA MECÁNICA", "3538 INGENIERÍA MECÁNICA", "3525 INGENIERÍA QUÍMICA", "3539 INGENIERÍA QUÍMICA"
]
const keyFacultad = "3068 FACULTAD DE MINAS";
const minasmerged = {};
keysMinas.forEach(key => {
    if (chunk2[keyFacultad].hasOwnProperty(key)) {
        minasmerged[key] = chunk2[keyFacultad][key];
    } else if (chunk3[keyFacultad].hasOwnProperty(key)) {
        minasmerged[key] = chunk3[keyFacultad][key];
    }
});

const merged = {
    ...chunk1,
    ...chunk4,
    "3068 FACULTAD DE MINAS": minasmerged,
}

fs.writeFileSync('data.json', JSON.stringify(merged));