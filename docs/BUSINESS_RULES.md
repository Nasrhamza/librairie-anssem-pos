# Librairie Anssem POS — Règles métier actuelles

## Vente / caisse
- 5 paniers clients indépendants et persistants.
- Paiement autorisé : **Espèces** ou **Crédit** uniquement.
- Pour une vente à crédit : nom client + téléphone (au moins 8 chiffres) obligatoires.
- Le stock physique est décrémenté uniquement après validation complète de la vente.
- Les services n'ont pas de stock.
- Le prix unitaire enregistré dans une vente reste celui du moment de la vente.
- Une ligne « Autre frais » sans produit est autorisée si sa description et son montant sont valides.
- Chaque vente peut être exportée comme facture client A4 en PDF.

## Produits
- Code-barres obligatoire et unique.
- Prix achat et prix vente modifiables.
- Un changement de prix vente ajoute une entrée à l'historique de prix.
- Scanner USB/HID : le scanner doit agir comme clavier et envoyer Enter après le code.
- Les produits peuvent être épinglés et sont regroupés par catégorie dans la caisse.

## Crédit client
- Un client est retrouvé par son numéro de téléphone normalisé.
- Une vente crédit augmente sa dette.
- Un paiement partiel ou total réduit la dette.
- Le logiciel refuse un paiement supérieur à la dette enregistrée.

## Connexion / SIM
- Chaque SIM source possède un quota journalier configurable (1000 Mo par défaut).
- Forfaits actuellement vendables : 500 Mo et 1 Go.
- Donc avec quota 1000 Mo : 1 x 1 Go OU 2 x 500 Mo au maximum sur la même SIM le même jour.
- Le logiciel bloque tout dépassement du quota journalier.
- Une SIM dont les 1000 Mo sont utilisés reste bloquée jusqu'au jour suivant.
- La suppression d'une SIM ne supprime jamais l'historique de ses ventes.
- Seuls les prix de vente 500 Mo / 1 Go sont exposés dans l'interface.
- Le rapport PDF journalier contient les détails et le total des ventes SIM.
- Le transfert réel opérateur n'est pas automatisé sans API/modem/USSD officiel.

## Données
- Sauvegarde locale automatique après les opérations métier.
- Écriture atomique : `data.json.tmp` puis renommage vers `data.json`.
- Windows : `%LOCALAPPDATA%\LibrairieAnssem\data.json`.
- Le bouton Backup exporte une copie JSON complète.

## Licence commerciale
- Aucune API métier n'est accessible avant activation valide.
- Chaque installation Windows génère un Code machine protégé par DPAPI au niveau du PC.
- Une Clé produit est signée Ed25519 par l'outil propriétaire et liée à un seul Code machine.
- Une clé copiée vers une autre installation est refusée.
- La licence est permanente pour le PC, sans connexion Internet obligatoire.
- Seul `dist\Jamel-v1-Setup.exe` peut être distribué au client.
- Le générateur et la clé privée dans `owner-tools` restent strictement propriétaires.
