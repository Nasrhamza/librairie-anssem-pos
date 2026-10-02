# Instructions pour Codex — Jamel v1 / Librairie Anssem POS

Ce dossier est la **source officielle qui correspond à l'EXE fourni**.

## Fichiers à connaître
- `main.go` : backend local, logique métier, stockage JSON, API HTTP locale.
- `web/index.html` : **interface exacte de l'EXE** (HTML/CSS/JS embarqué).
- `desktop_windows.go` : fenêtre Windows native WebView2 et icône de l'application.
- `assets/jamel-v1.ico` + `app_icon_windows_amd64.syso` : icône Windows embarquée.
- `main_test.go` : tests de base.
- `docs/BUSINESS_RULES.md` : règles métier à préserver.
- `build_windows.bat` : tests + compilation EXE Windows.
- `installer/main.go` : installeur Windows autonome, raccourcis et désinstallation.
- `installer/payload/MicrosoftEdgeWebview2Setup.exe` : bootstrapper officiel Microsoft utilisé seulement si WebView2 manque sur le PC client.
- `license_manager.go` + `license_windows.go` : activation obligatoire, signature de licence et liaison DPAPI au PC Windows.
- `internal/licensing` : format et vérification cryptographique des clés produit.
- `owner-tools/Jamel-License-Generator.exe` : outil privé de création des clés, à ne jamais distribuer.

## Règle principale
Ne jamais créer une deuxième interface séparée de `web/index.html`.
Toute modification visuelle destinée au logiciel doit être faite dans `web/index.html`, puis reconstruire l'EXE avec `build_windows.bat`.

## Contraintes métier à ne pas casser
1. Cinq paniers clients indépendants doivent survivre à une fermeture/redémarrage.
2. Paiement : Espèces ou Crédit seulement.
3. Crédit => nom + téléphone obligatoires.
4. Vente stock => vérifier tout le panier avant de décrémenter le stock.
5. Le code-barres doit rester unique.
6. Les anciennes ventes gardent leur ancien prix.
7. SIM : quota journalier appliqué avant enregistrement d'une vente connexion.
8. Les données doivent rester locales et persistantes.
9. Ne pas simuler un transfert télécom réel si aucun mécanisme opérateur réel n'est configuré.
10. Ne jamais retirer ni contourner le contrôle de licence avant les API métier.
11. Ne jamais copier la clé privée de `owner-tools/license-generator/private` dans le Setup, le dossier `dist` ou une réponse destinée au client.

## Après modification
Exécuter :

```bat
build_windows.bat
```

Le nouveau programme sera :
`dist\Jamel-v1.exe`

Le fichier à envoyer à un autre PC sera :
`dist\Jamel-v1-Setup.exe`

Le générateur propriétaire sera reconstruit séparément dans :
`owner-tools\Jamel-License-Generator.exe`
