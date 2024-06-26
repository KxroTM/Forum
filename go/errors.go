package forum

import "errors"

// TOUTES LES ERREURS COMPORTES DES MAJ CAR ELLES SONT UTILISEES POUR LES AFFICHER DANS LES TEMPLATES

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
var ErrSpaceInUsername = errors.New("Le pseudo ne doit pas contenir d'espace ou de '-', si vous voulez en utilisez un, utilisez un '_' à la place")
var ErrSpaceInPassword = errors.New("Le mot de passe ne doit pas contenir d'espace")
var ErrBadTypeImg = errors.New("Type d'image invalide")
var ErrBadSizeImg = errors.New("Taille de l'image invalide. La taille maximale est de 20Mo")
var ErrEmailAlreadyUsed = errors.New("Email déjà utilisé")
var ErrNotPremium = errors.New("Vous n'êtes pas premium, vous ne pouvez pas utilisé d'image GIF")
