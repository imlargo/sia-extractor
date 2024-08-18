const DATA = require("./data.json");
const { MongoClient, ServerApiVersion } = require('mongodb');

const uri = process.env.MONGO_URI;

const client = new MongoClient(uri, {
	serverApi: {
		version: ServerApiVersion.v1,
		strict: true,
		deprecationErrors: true,
	}
});


(async () => {
	
	await client.connect();

	const collAsignaturas = client.db("asignaturas").collection("asignaturas");

	const promises = Object.entries(DATA).map(async ([facultad, data]) => {

		const query = { _id: facultad };
		const result = await collAsignaturas.replaceOne(
			query, data
		);

		console.log(`Facultad ${facultad} actualizada`);

		return result.modifiedCount
	});

	const results = await Promise.all(promises);

	if (results.some((result) => result === 0)){
		console.log("Some documents were not updated");
		throw new Error("Some documents were not updated");
	}

	console.log("Datos actualizados con exito");

	await client.close();
})();