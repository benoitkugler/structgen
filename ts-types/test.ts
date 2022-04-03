// Code generated by structgen DO NOT EDIT
// github.com/benoitkugler/goACVE/server/core/rawdata.Euros
export type Euros = number;
// github.com/benoitkugler/goACVE/server/core/rawdata.Bool
export type Bool = boolean;
// github.com/benoitkugler/goACVE/server/core/rawdata.Int
export type Int = number;
// github.com/benoitkugler/goACVE/server/core/rawdata.Aide
export interface Aide {
  id: number;
  id_structureaide: number;
  id_participant: number;
  valeur: Euros;
  valide: Bool;
  par_jour: Bool;
  nb_jours_max: Int;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.String
export type String_ = string;
// github.com/benoitkugler/goACVE/server/core/rawdata.TrajetBus
export interface TrajetBus {
  rendez_vous: String_;
  prix: Euros;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.BusCamp
export interface BusCamp {
  actif: boolean;
  commentaire: String_;
  aller: TrajetBus;
  retour: TrajetBus;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MaterielSkiCamp
export interface MaterielSkiCamp {
  actif: boolean;
  prix_acve: Euros;
  prix_loueur: Euros;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionsCamp
export interface OptionsCamp {
  bus: BusCamp;
  materiel_ski: MaterielSkiCamp;
}

class DateTag {
  private _: "D" = "D";
}

class TimeTag {
  private _: "T" = "T";
}

// AAAA-MM-YY date format
export type Date_ = string & DateTag;

// ISO date-time string
export type Time = string & TimeTag;

// github.com/benoitkugler/goACVE/server/core/rawdata.Vetement
export interface Vetement {
  quantite: number;
  description: string;
  obligatoire: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ListeVetements
export interface ListeVetements {
  liste: Vetement[] | null;
  complement: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.SchemaPaiement
enum SchemaPaiement {
  SPAcompte = "acompte",
  SPTotal = "total"
}

export const SchemaPaiementLabels: { [key in SchemaPaiement]: string } = {
  [SchemaPaiement.SPAcompte]: "Avec acompte",
  [SchemaPaiement.SPTotal]: "Paiement direct (sans acompte)"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Envois
export interface Envois {
  __locked__: boolean;
  lettre_directeur: boolean;
  liste_vetements: boolean;
  liste_participants: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Plage
export interface Plage {
  from: Date_;
  to: Date_;
}

// github.com/benoitkugler/goACVE/server/core/rawdata.OptionSemaineCamp
export interface OptionSemaineCamp {
  plage_1: Plage;
  plage_2: Plage;
  prix_1: Euros;
  prix_2: Euros;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.PrixParStatut
export interface PrixParStatut {
  id: number;
  prix: Euros;
  statut: String_;
  description: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionPrixCamp
export interface OptionPrixCamp {
  active: string;
  semaine: OptionSemaineCamp;
  statut: PrixParStatut[] | null;
  jour: Euros[] | null;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionQuotientFamilial
export type OptionQuotientFamilial = number[];
// github.com/benoitkugler/goACVE/server/core/rawdata.Camp
export interface Camp {
  id: number;
  lieu: String_;
  nom: String_;
  prix: Euros;
  nb_places: Int;
  password: String_;
  ouvert: Bool;
  nb_places_reservees: Int;
  numero_js: String_;
  need_equilibre_gf: Bool;
  age_min: Int;
  age_max: Int;
  options: OptionsCamp;
  date_debut: Date_;
  date_fin: Date_;
  liste_vetements: ListeVetements;
  schema_paiement: SchemaPaiement;
  joomeo_album_id: String_;
  envois: Envois;
  lien_compta: String_;
  option_prix: OptionPrixCamp;
  inscription_simple: Bool;
  infos: String_;
  quotient_familial: OptionQuotientFamilial;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.CampContrainte
export interface CampContrainte {
  id_camp: number;
  id_contrainte: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ContenuDocument
export interface ContenuDocument {
  id_document: number;
  contenu: number[] | null;
  miniature: number[] | null;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionnalId
export interface OptionnalId {
  Int64: number;
  Valid: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.BuiltinContrainte
enum BuiltinContrainte {
  CAutre = "autre",
  CBafa = "bafa",
  CBafaEquiv = "bafa_equiv",
  CBafd = "bafd",
  CBafdEquiv = "bafd_equiv",
  CCarteId = "carte_id",
  CCarteVitale = "carte_vitale",
  CCertMedCuisine = "cert_med_cuisine",
  CHaccp = "haccp",
  CInvalide = "",
  CPermis = "permis",
  CSb = "sb",
  CScolarite = "scolarite",
  CSecour = "secour",
  CTestNautique = "test_nautique",
  CVaccin = "vaccin"
}

export const BuiltinContrainteLabels: { [key in BuiltinContrainte]: string } = {
  [BuiltinContrainte.CAutre]: "Autre",
  [BuiltinContrainte.CBafa]: "BAFA",
  [BuiltinContrainte.CBafaEquiv]: "Equivalent BAFA",
  [BuiltinContrainte.CBafd]: "BAFD",
  [BuiltinContrainte.CBafdEquiv]: "Equivalent BAFD",
  [BuiltinContrainte.CCarteId]: "Carte d''identité/Passeport",
  [BuiltinContrainte.CCarteVitale]: "Carte Vitale",
  [BuiltinContrainte.CCertMedCuisine]: "Certificat médical Cuisine",
  [BuiltinContrainte.CHaccp]: "Cuisine (HACCP)",
  [BuiltinContrainte.CInvalide]: "-",
  [BuiltinContrainte.CPermis]: "Permis de conduire",
  [BuiltinContrainte.CSb]: "Surveillant de baignade",
  [BuiltinContrainte.CScolarite]: "Certificat de scolarité",
  [BuiltinContrainte.CSecour]: "Secourisme (PSC1 - AFPS)",
  [BuiltinContrainte.CTestNautique]: "Test nautique",
  [BuiltinContrainte.CVaccin]: "Vaccin"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Contrainte
export interface Contrainte {
  id: number;
  id_personne: OptionnalId;
  id_document: OptionnalId;
  builtin: BuiltinContrainte;
  nom: String_;
  description: String_;
  max_docs: number;
  jours_valide: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Taille
export type Taille = number;
// github.com/benoitkugler/goACVE/server/core/rawdata.Document
export interface Document {
  id: number;
  taille: Taille;
  nom_client: String_;
  description: String_;
  date_heure_modif: Time;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.DocumentAide
export interface DocumentAide {
  id_document: number;
  id_aide: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.DocumentCamp
export interface DocumentCamp {
  id_document: number;
  id_camp: number;
  is_lettre: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.DocumentPersonne
export interface DocumentPersonne {
  id_document: number;
  id_personne: number;
  id_contrainte: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ModePaiment
enum ModePaiment {
  MPAncv = "ancv",
  MPAucun = "",
  MPCarte = "cb",
  MPCheque = "cheque",
  MPEspece = "esp",
  MPHelloasso = "helloasso",
  MPVirement = "vir"
}

export const ModePaimentLabels: { [key in ModePaiment]: string } = {
  [ModePaiment.MPAncv]: "ANCV",
  [ModePaiment.MPAucun]: "-",
  [ModePaiment.MPCarte]: "Carte bancaire",
  [ModePaiment.MPCheque]: "Chèque",
  [ModePaiment.MPEspece]: "Espèces",
  [ModePaiment.MPHelloasso]: "Hello Asso",
  [ModePaiment.MPVirement]: "Virement"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.InfoDon
export interface InfoDon {
  id_paiement_hello_asso: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Don
export interface Don {
  id: number;
  valeur: Euros;
  mode_paiement: ModePaiment;
  date_reception: Date_;
  recu_emis: Date_;
  infos: InfoDon;
  remercie: Bool;
  details: String_;
  affectation: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.DonDonateur
export interface DonDonateur {
  id_don: number;
  id_personne: OptionnalId;
  id_organisme: OptionnalId;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Role
enum Role {
  RAdjoint = "_adjoint",
  RAideAnimation = "_aideanim",
  RAnimation = "_anim",
  RAutre = "_autre",
  RBabysiter = "_babysiter",
  RChauffeur = "_chauffeur",
  RCuis = "_cuis",
  RDirecteur = "_dir",
  RFactotum = "_factotum",
  RInfirm = "_infirm",
  RIntend = "_intend",
  RLing = "_ling",
  RMen = "_men"
}

export const RoleLabels: { [key in Role]: string } = {
  [Role.RAdjoint]: "Adjoint",
  [Role.RAideAnimation]: "Aide-animateur",
  [Role.RAnimation]: "Animation",
  [Role.RAutre]: "Autre",
  [Role.RBabysiter]: "Baby-sitter",
  [Role.RChauffeur]: "Chauffeur",
  [Role.RCuis]: "Cuisine",
  [Role.RDirecteur]: "Direction",
  [Role.RFactotum]: "Factotum",
  [Role.RInfirm]: "Assistant sanitaire",
  [Role.RIntend]: "Intendance",
  [Role.RLing]: "Lingerie",
  [Role.RMen]: "Ménage"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Roles
export type Roles = Role[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.Diplome
enum Diplome {
  DAgreg = "agreg",
  DAssSociale = "ass_sociale",
  DAucun = "",
  DBafa = "bafa",
  DBafaStag = "bafa_stag",
  DBafd = "bafd",
  DBafdStag = "bafd_stag",
  DBapaat = "bapaat",
  DBeatep = "beatep",
  DBjeps = "bjeps",
  DCap = "cap",
  DDeug = "deug",
  DDut = "dut",
  DEducSpe = "educ_spe",
  DEje = "eje",
  DInstit = "instit",
  DMonEduc = "mon_educ",
  DProf = "prof",
  DStaps = "staps",
  DZzautre = "zzautre"
}

export const DiplomeLabels: { [key in Diplome]: string } = {
  [Diplome.DAgreg]: "Agrégé",
  [Diplome.DAssSociale]: "Assitante Sociale",
  [Diplome.DAucun]: "Aucun",
  [Diplome.DBafa]: "BAFA Titulaire",
  [Diplome.DBafaStag]: "BAFA Stagiaire",
  [Diplome.DBafd]: "BAFD titulaire",
  [Diplome.DBafdStag]: "BAFD stagiaire",
  [Diplome.DBapaat]: "BAPAAT",
  [Diplome.DBeatep]: "BEATEP",
  [Diplome.DBjeps]: "BPJEPS",
  [Diplome.DCap]: "CAP petit enfance",
  [Diplome.DDeug]: "DEUG",
  [Diplome.DDut]: "DUT carrière sociale",
  [Diplome.DEducSpe]: "Educ. spé.",
  [Diplome.DEje]: "EJE",
  [Diplome.DInstit]: "Professeur des écoles",
  [Diplome.DMonEduc]: "Moniteur educateur",
  [Diplome.DProf]: "Enseignant du secondaire",
  [Diplome.DStaps]: "STAPS",
  [Diplome.DZzautre]: "AUTRE"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Approfondissement
enum Approfondissement {
  AAucun = "",
  AAutre = "autre",
  ACanoe = "canoe",
  AMoto = "moto",
  ASb = "sb",
  AVoile = "voile"
}

export const ApprofondissementLabels: { [key in Approfondissement]: string } = {
  [Approfondissement.AAucun]: "Non effectué",
  [Approfondissement.AAutre]: "Approfondissement",
  [Approfondissement.ACanoe]: "Canoë - Kayak",
  [Approfondissement.AMoto]: "Loisirs motocyclistes",
  [Approfondissement.ASb]: "Surveillant de baignade",
  [Approfondissement.AVoile]: "Voile"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.OptionnalPlage
export type OptionnalPlage = {
  active: boolean;
} & Plage;
// github.com/benoitkugler/goACVE/server/core/rawdata.InvitationEquipier
export type InvitationEquipier = number;
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionnalBool
enum OptionnalBool {
  OBNon = -1,
  OBOui = 1,
  OBPeutEtre = 0
}

export const OptionnalBoolLabels: { [key in OptionnalBool]: string } = {
  [OptionnalBool.OBNon]: "Non",
  [OptionnalBool.OBOui]: "Oui",
  [OptionnalBool.OBPeutEtre]: "Peut-être"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Equipier
export interface Equipier {
  id: number;
  id_camp: number;
  id_personne: number;
  roles: Roles;
  diplome: Diplome;
  appro: Approfondissement;
  presence: OptionnalPlage;
  invitation_equipier: InvitationEquipier;
  charte: OptionnalBool;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.EquipierContrainte
export interface EquipierContrainte {
  id_equipier: number;
  id_contrainte: number;
  optionnel: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Sexe
enum Sexe {
  SAucun = "",
  SFemme = "F",
  SHomme = "M"
}

export const SexeLabels: { [key in Sexe]: string } = {
  [Sexe.SAucun]: "-",
  [Sexe.SFemme]: "Femme",
  [Sexe.SHomme]: "Homme"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Destinataire
export interface Destinataire {
  nom_prenom: String_;
  sexe: Sexe;
  adresse: String_;
  code_postal: String_;
  ville: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.DestinatairesOptionnels
export type DestinatairesOptionnels = Destinataire[] | null;
// github.com/lib/pq.StringArray
export type StringArray = string[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.Facture
export interface Facture {
  id: number;
  id_personne: number;
  destinataires_optionnels: DestinatairesOptionnels;
  key: String_;
  copies_mails: StringArray;
  last_connection: Time;
  is_confirmed: boolean;
  is_validated: boolean;
  partage_adresses_ok: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Groupe
export interface Groupe {
  id: number;
  id_camp: number;
  nom: String_;
  plage: Plage;
  couleur: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.GroupeContrainte
export interface GroupeContrainte {
  id_groupe: number;
  id_contrainte: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.GroupeEquipier
export interface GroupeEquipier {
  id_groupe: number;
  id_equipier: number;
  id_camp: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.GroupeParticipant
export interface GroupeParticipant {
  id_participant: number;
  id_groupe: number;
  id_camp: number;
  manuel: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Imageuploaded
export interface Imageuploaded {
  id_camp: number;
  filename: string;
  lien: string;
  content: number[] | null;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.IdentificationId
export interface IdentificationId {
  valid: boolean;
  id: number;
  crypted: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Tels
export type Tels = string[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.Pays
export type Pays = string;
// github.com/benoitkugler/goACVE/server/core/rawdata.ResponsableLegal
export interface ResponsableLegal {
  lienid: IdentificationId;
  nom: String_;
  prenom: String_;
  sexe: Sexe;
  mail: String_;
  adresse: String_;
  code_postal: String_;
  ville: String_;
  tels: Tels;
  date_naissance: Date_;
  pays: Pays;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Bus
enum Bus {
  BAller = "aller",
  BAllerRetour = "aller_retour",
  BAucun = "",
  BRetour = "retour"
}

export const BusLabels: { [key in Bus]: string } = {
  [Bus.BAller]: "Aller",
  [Bus.BAllerRetour]: "Aller-Retour",
  [Bus.BAucun]: "-",
  [Bus.BRetour]: "Retour"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.MaterielSki
export interface MaterielSki {
  need: string;
  mode: string;
  casque: boolean;
  poids: number;
  taille: number;
  pointure: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionsParticipant
export interface OptionsParticipant {
  bus: Bus;
  materiel_ski: MaterielSki;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Semaine
enum Semaine {
  SComplet = "",
  SSe1 = "1",
  SSe2 = "2"
}

export const SemaineLabels: { [key in Semaine]: string } = {
  [Semaine.SComplet]: "Camp complet",
  [Semaine.SSe1]: "Semaine 1",
  [Semaine.SSe2]: "Semaine 2"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Jours
export type Jours = number[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.OptionPrixParticipant
export interface OptionPrixParticipant {
  semaine: Semaine;
  statut: number;
  jour: Jours;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ParticipantInscription
export interface ParticipantInscription {
  lienid: IdentificationId;
  nom: String_;
  prenom: String_;
  date_naissance: Date_;
  sexe: Sexe;
  id_camp: number;
  options: OptionsParticipant;
  options_prix: OptionPrixParticipant;
  quotient_familial: Int;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ParticipantInscriptions
export type ParticipantInscriptions = ParticipantInscription[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.Inscription
export interface Inscription {
  id: number;
  info: String_;
  date_heure: Time;
  copies_mails: StringArray;
  responsable: ResponsableLegal;
  participants: ParticipantInscriptions;
  partage_adresses_ok: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Lettredirecteur
export interface Lettredirecteur {
  id_camp: number;
  html: string;
  use_coord_centre: boolean;
  show_adresse_postale: boolean;
  color_coord: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessageKind
enum MessageKind {
  MAccuseReception = 3,
  MAttestationPresence = 7,
  MCentre = 2,
  MDocuments = 5,
  MFacture = 4,
  MFactureAcquittee = 6,
  MInscription = 9,
  MPaiement = 11,
  MPlaceLiberee = 10,
  MResponsable = 1,
  MSondage = 8,
  MSupprime = 0
}

export const MessageKindLabels: { [key in MessageKind]: string } = {
  [MessageKind.MAccuseReception]: "Inscription validée",
  [MessageKind.MAttestationPresence]: "Attestation de présence",
  [MessageKind.MCentre]: "Message du centre",
  [MessageKind.MDocuments]: "Document des séjours",
  [MessageKind.MFacture]: "Facture",
  [MessageKind.MFactureAcquittee]: "Facture acquittée",
  [MessageKind.MInscription]: "Moment d'inscription",
  [MessageKind.MPaiement]: "",
  [MessageKind.MPlaceLiberee]: "Place libérée",
  [MessageKind.MResponsable]: "Message",
  [MessageKind.MSondage]: "Avis sur le séjour",
  [MessageKind.MSupprime]: "Message supprimé"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Message
export interface Message {
  id: number;
  id_facture: number;
  kind: MessageKind;
  created: Time;
  modified: Time;
  vu: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Distribution
enum Distribution {
  DEspacePerso = 0,
  DMail = 1,
  DMailAndDownload = 2
}

export const DistributionLabels: { [key in Distribution]: string } = {
  [Distribution.DEspacePerso]: "Téléchargée depuis l'espace de suivi",
  [Distribution.DMail]: "Notifiée par courriel",
  [Distribution.DMailAndDownload]: "Téléchargée après notification"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.MessageAttestation
export interface MessageAttestation {
  id_message: number;
  distribution: Distribution;
  guard_kind: MessageKind;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessageDocument
export interface MessageDocument {
  id_message: number;
  id_camp: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessageMessage
export interface MessageMessage {
  id_message: number;
  contenu: String_;
  guard_kind: MessageKind;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessagePlacelibere
export interface MessagePlacelibere {
  id_message: number;
  id_participant: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessageSondage
export interface MessageSondage {
  id_message: number;
  id_camp: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.MessageView
export interface MessageView {
  id_message: number;
  id_camp: number;
  vu: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Coordonnees
export interface Coordonnees {
  tels: Tels;
  mail: String_;
  adresse: String_;
  code_postal: String_;
  ville: String_;
  pays: Pays;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Exemplaires
export interface Exemplaires {
  pub_ete: number;
  pub_hiver: number;
  echo_rocher: number;
  e_onews: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Organisme
export interface Organisme {
  id: number;
  nom: String_;
  contact_propre: Bool;
  contact: Coordonnees;
  id_contact: OptionnalId;
  id_contact_don: OptionnalId;
  exemplaires: Exemplaires;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Paiement
export interface Paiement {
  id: number;
  id_facture: number;
  is_acompte: Bool;
  is_remboursement: Bool;
  in_bordereau: Time;
  label_payeur: String_;
  nom_banque: String_;
  mode_paiement: ModePaiment;
  numero: String_;
  valeur: Euros;
  is_invalide: Bool;
  date_reglement: Time;
  details: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.StatutAttente
enum StatutAttente {
  Attente = 1,
  AttenteReponse = 2,
  Inscrit = 0,
  Refuse = 3
}

export const StatutAttenteLabels: { [key in StatutAttente]: string } = {
  [StatutAttente.Attente]: "Liste d'attente",
  [StatutAttente.AttenteReponse]: "Attente de confirmation",
  [StatutAttente.Inscrit]: "Inscrit",
  [StatutAttente.Refuse]: "Refusé"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.ListeAttente
export interface ListeAttente {
  statut: StatutAttente;
  raison: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Pourcent
export type Pourcent = number;
// github.com/benoitkugler/goACVE/server/core/rawdata.Remises
export interface Remises {
  reduc_equipiers: Pourcent;
  reduc_enfants: Pourcent;
  reduc_speciale: Euros;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Participant
export interface Participant {
  id: number;
  id_camp: number;
  id_personne: number;
  id_facture: OptionnalId;
  liste_attente: ListeAttente;
  remises: Remises;
  option_prix: OptionPrixParticipant;
  options: OptionsParticipant;
  date_heure: Time;
  quotient_familial: Int;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.ParticipantEquipier
export interface ParticipantEquipier {
  id_participant: number;
  id_equipier: number;
  id_groupe: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Participantsimple
export interface Participantsimple {
  id: number;
  id_personne: number;
  id_camp: number;
  date_heure: Time;
  info: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.RangMembreAsso
enum RangMembreAsso {
  RMABureau = "3",
  RMACA = "2",
  RMAMembre = "1",
  RMANonMembre = ""
}

export const RangMembreAssoLabels: { [key in RangMembreAsso]: string } = {
  [RangMembreAsso.RMABureau]: "Membre du bureau",
  [RangMembreAsso.RMACA]: "Membre du C.A.",
  [RangMembreAsso.RMAMembre]: "Membre",
  [RangMembreAsso.RMANonMembre]: "Non membre"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.Cotisation
export type Cotisation = number[] | null;
// github.com/benoitkugler/goACVE/server/core/rawdata.Maladies
export interface Maladies {
  rubeole: boolean;
  varicelle: boolean;
  angine: boolean;
  oreillons: boolean;
  scarlatine: boolean;
  coqueluche: boolean;
  otite: boolean;
  rougeole: boolean;
  rhumatisme: boolean;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Allergies
export interface Allergies {
  asthme: boolean;
  alimentaires: boolean;
  medicamenteuses: boolean;
  autres: string;
  conduite_a_tenir: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Medecin
export interface Medecin {
  nom: string;
  tel: string;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.FicheSanitaire
export interface FicheSanitaire {
  traitement_medical: boolean;
  maladies: Maladies;
  allergies: Allergies;
  difficultes_sante: string;
  recommandations: string;
  handicap: boolean;
  tel: string;
  medecin: Medecin;
  last_modif: Time;
  mails: string[] | null;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Departement
export type Departement = string;
// github.com/benoitkugler/goACVE/server/core/rawdata.BasePersonne
export interface BasePersonne {
  nom: String_;
  nom_jeune_fille: String_;
  prenom: String_;
  date_naissance: Date_;
  ville_naissance: String_;
  departement_naissance: Departement;
  sexe: Sexe;
  tels: Tels;
  mail: String_;
  adresse: String_;
  code_postal: String_;
  ville: String_;
  pays: Pays;
  securite_sociale: String_;
  profession: String_;
  etudiant: Bool;
  fonctionnaire: Bool;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Personne
export type Personne = {
  id: number;
  version_papier: Bool;
  pub_hiver: Bool;
  pub_ete: Bool;
  echo_rocher: Bool;
  rang_membre_asso: RangMembreAsso;
  quotient_familial: Int;
  cotisation: Cotisation;
  eonews: Bool;
  fiche_sanitaire: FicheSanitaire;
  is_temporaire: Bool;
} & BasePersonne;
// github.com/benoitkugler/goACVE/server/core/rawdata.Satisfaction
enum Satisfaction {
  SDecevant = 1,
  SMoyen = 2,
  SSatisfaisant = 3,
  STressatisfaisant = 4,
  SVide = 0
}

export const SatisfactionLabels: { [key in Satisfaction]: string } = {
  [Satisfaction.SDecevant]: "Décevant",
  [Satisfaction.SMoyen]: "Moyen",
  [Satisfaction.SSatisfaisant]: "Satisfaisant",
  [Satisfaction.STressatisfaisant]: "Très satisfaisant",
  [Satisfaction.SVide]: "-"
};

// github.com/benoitkugler/goACVE/server/core/rawdata.RepSondage
export interface RepSondage {
  infos_avant_sejour: Satisfaction;
  infos_pendant_sejour: Satisfaction;
  hebergement: Satisfaction;
  activites: Satisfaction;
  theme: Satisfaction;
  nourriture: Satisfaction;
  hygiene: Satisfaction;
  ambiance: Satisfaction;
  ressenti: Satisfaction;
  message_enfant: String_;
  message_responsable: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Sondage
export type Sondage = {
  id: number;
  id_camp: number;
  id_facture: number;
  modified: Time;
} & RepSondage;
// github.com/benoitkugler/goACVE/server/core/rawdata.Structureaide
export interface Structureaide {
  id: number;
  nom: String_;
  immatriculation: String_;
  adresse: String_;
  code_postal: String_;
  ville: String_;
  telephone: String_;
  info: String_;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.Modules
export interface Modules {
  personnes: number;
  camps: number;
  inscriptions: number;
  suivi_camps: number;
  suivi_dossiers: number;
  paiements: number;
  aides: number;
  equipiers: number;
  dons: number;
}
// github.com/benoitkugler/goACVE/server/core/rawdata.User
export interface User {
  id: number;
  label: String_;
  mdp: String_;
  is_admin: Bool;
  modules: Modules;
}

interface Itf {
  Kind: number;
  Data: User | Modules;
}
