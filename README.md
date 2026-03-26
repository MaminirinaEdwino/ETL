# ETL CLI

ETL CLI est un outils ETL en ligne de commande simple  permettant d'extraire de transformer et de charger les données souhaités par l'utilisateur dans un fichier JSON à partir d'un source de donnée en CSV contenant des millier de ligne.

## Note

La version actuel de l'outils ne permet pas encore d'inserer les données dans une base de données mais cela viendra dans la version future de l'outils.

## Utilisation

Pour lancer l'outils il faut utilisé la commande suivante :

```bash
etl --inputfile="road_accident_data.csv" --outputfile="res.json"
```

Après le lancement, un ecran affichant la liste des champs dans le source de donnée apparait. Pour choisir les champs a extraire on utilise les touches directionnel haut et bas pour naviguer entre les choix et entré pour les selectionnée.

Et pour ajouter des filtres comme par exemple nombre de voiture ou le jour de l'accident pour notre exemple , il faut tapez la touche `f` du clavier et passé a l'cran suivant pour choisir le filtre.

(Pour naviguer entre les onglets du cli il faut utilisé la touche `ctrl+gauche` ou `ctrl+droite`)

Dans le deuxième onglet, on utilise la touche `tab` pour mettre le focus sur le champs de saisie et pour saisir la valer souhaité dans le fitlre.
Puis pour ajouter le type du filtre, on appuie sur la touche `t` et utilisé les touches directionnel pour se positionner sur le choix, Pour le choix de l'operation il faut taper sur la touche `o` et executé les meme action qu'avec le type. et enfin pour ajouter les valuers du filtre on appuie sur f et le filtre appraitra sur le bas de l'onglet des filtres.

Et pour terminer avec le chargement des données dans le fichier json, on passe au troisième onglet et  on appuie sur `ctrl+e` pour les charger et afficher un éxtrait des résultat .
