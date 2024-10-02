parallel -j 20 --ungroup 'go run . extract {1}' ::: {1..20}
# MONGO_URI="mongodb+srv://imlargo:VQAWP8qMxhhsp3aD@asignaturas.vs2qizj.mongodb.net/?retryWrites=true&w=majority&appName=asignaturas" go run . deploy
