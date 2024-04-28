const fs = require('fs/promises');

async function loadJson(path) {
    const data = await fs.readFile(path, 'utf-8');
    return JSON.parse(data);
}

async function main() {

    const indexs = Array.from({ length: 9 }, (_, i) => i + 1);

    const data = await Promise.all(
        indexs.map(async i => await loadJson(`../artifacts/${i}.json`))
    );

    await fs.writeFile('data.json', JSON.stringify(data));
}

main();