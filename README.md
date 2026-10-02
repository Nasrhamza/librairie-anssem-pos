# Jamel v1 — Librairie Anssem POS

Ce dossier contient **le source et l'EXE issus de la même interface**.
Il n'y a plus une version visuelle séparée pour le source et une autre pour l'EXE.

## Installer sur ce PC ou sur un autre PC
Copier puis ouvrir :
`dist\Jamel-v1-Setup.exe`

Le Setup fonctionne sur Windows 10/11 64 bits sans installer Go. Il installe le
logiciel dans le profil Windows, crée les raccourcis Bureau et menu Démarrer et
ajoute une désinstallation Windows. Le logiciel s'ouvre dans sa propre fenêtre
desktop WebView2, sans onglet ni barre d'adresse Edge/Chrome. Le Setup contient
le bootstrapper officiel Microsoft et installe WebView2 silencieusement si le
PC client ne le possède pas encore.
Le logo livre/crayon de Librairie Anssem est intégré comme icône de la fenêtre,
du fichier EXE, du Setup et des raccourcis Windows.

Au premier démarrage, le propriétaire crée son nom administrateur et son mot de
passe. Il n'existe aucun mot de passe par défaut.

Pour tester la version portable sans installation :
`dist\Jamel-v1.exe`

## Activation commerciale par clé produit
Jamel v1 est bloqué avant la connexion administrateur tant qu'il n'est pas
activé. Chaque nouvelle installation affiche un **Code machine**. Le
propriétaire crée ensuite une **Clé produit unique** avec :
`owner-tools\Jamel-License-Generator.exe`

La clé produit est signée numériquement et liée à l'installation Windows. Une
clé créée pour un PC est refusée sur un autre PC. Le système fonctionne hors
ligne : le client peut envoyer son Code machine par WhatsApp, puis coller la clé
reçue dans l'écran d'activation.

Ne jamais envoyer au client le dossier `owner-tools`, le générateur, la clé
privée ou le dossier source. Le seul fichier à distribuer est :
`dist\Jamel-v1-Setup.exe`

La procédure propriétaire complète se trouve dans :
`owner-tools\NE-JAMAIS-ENVOYER-AU-CLIENT.txt`

## Données persistantes
Sous Windows :
`%LOCALAPPDATA%\LibrairieAnssem\data.json`

Le compte administrateur est protégé séparément dans :
`%LOCALAPPDATA%\LibrairieAnssem\auth.json`

La licence et la liaison Windows de la machine sont conservées dans le même
dossier. Copier la licence vers un autre ordinateur ne l'active pas.

Lors d'une mise à jour depuis un ancien nom du logiciel, les données et le
compte administrateur sont copiés automatiquement vers le nouveau dossier.

La fermeture du logiciel ou du PC ne supprime pas les données enregistrées.

## Modifier avec VS Code / Codex
- Interface : `web\index.html`
- Logique : `main.go`
- Règles : `docs\BUSINESS_RULES.md`
- Guide Codex : `CODEX_README.md`

## Reconstruire l'EXE
Prérequis : Go installé.
Double-cliquer :
`build_windows.bat`

Le script exécute les tests avant de compiler.

## Fonctions présentes dans cette version d'essai
- Dashboard.
- Caisse avec 5 clients/paniers simultanés.
- Scan code-barres USB/HID à la caisse.
- Produits en grille compacte continue, filtrables par catégorie et épinglables en tête de caisse.
- Ligne « Autre frais » pour les ventes sans produit ni code-barres.
- Produits et services.
- Ajout/modification des produits.
- Prix achat/vente modifiables + historique de prix vente.
- Stock et entrée/sortie manuelle.
- Vente Espèces.
- Vente Crédit avec nom + téléphone obligatoires.
- Clients débiteurs + encaissement partiel/total.
- SIM sources, quota journalier.
- Vente connexion 500 Mo / 1 Go depuis chaque ligne SIM.
- Ajout et suppression des SIM, avec conservation de l'historique.
- Prix connexion modifiables.
- Rapport PDF journalier des ventes SIM.
- Facture client A4 en PDF pour chaque vente.
- Facture personnalisée avec le logo Librairie Anssem, WhatsApp, email et matricule fiscal.
- Accès direct à une sélection de services digitaux officiels.
- Historique ventes et connexions.
- Backup JSON.
- Stockage local offline.
- Connexion administrateur au démarrage.
- Activation commerciale hors ligne par clé produit liée à un seul PC.
- Thème personnalisable avec sélecteur de couleur.
- Fenêtre desktop WebView2 native, indépendante d'Edge, et instance unique.
- Icône Windows officielle Librairie Anssem sur l'application et le Setup.
- Setup Windows avec raccourcis et désinstallation (les données sont conservées).

## Limite assumée
Le logiciel enregistre et contrôle une vente de connexion mais ne déclenche pas encore le transfert réel chez l'opérateur. Cela nécessitera la méthode réelle : API opérateur, modem GSM/USSD ou autre interface autorisée.
