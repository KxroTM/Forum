package forum

import "errors"

// TOUTES LES ERREURS COMPORTES DES MAJ CAR ELLES SONT UTILISEES DANS POUR LES AFFICHER DANS LES TEMPLATES

// DONNEES VIDE
var ErrEmptyFieldEmail = errors.New("Email non renseigné")
var ErrEmptyFieldPassword = errors.New("Mot de passe non renseigné")
var ErrEmptyFieldPasswordCheck = errors.New("Mot de passe de vérification non renseigné")
var ErrEmptyFieldPseudo = errors.New("Pseudo non renseigné")

// DONNEES INVALIDES
var ErrInvalidPasswordCheck = errors.New("Les mots de passe ne correspondent pas")
var ErrBadEmail = errors.New("Compte inexistant")
var ErrBadPassword = errors.New("Mot de passe incorrect")
var ErrInvalidPseudo = errors.New("Pseudo invalide. \n Le pseudo doit contenir entre 4 et 15 caractères et ne doit pas contenir des mots interdits")
var ErrInvalidEmail = errors.New("Email invalide")
var ErrMailAlreadyUsed = errors.New("Compte déjà existant")
var ErrPseudoAlreadyUsed = errors.New("Pseudo déjà utilisé")
var ErrInvalidPassword = errors.New("Mot de passe invalide. \n Le mot de passe doit contenir au moins 8 caractères, une majuscule, une minuscule, un chiffre et un caractère spécial")
