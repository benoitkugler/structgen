import 'dart:convert';

import 'test.dart';

void main() {
  final cl = Livraison(0, 012, "dsd", [true, false], 4, 5);

  final js = LivraisonToJson(cl);

  print("${js}");
  final cl2 = LivraisonFromJson(js);

  print("${cl2.anticipation}, ${cl2.jours_livraison}");

  final cm = Commande(2, 4, DateTime.now(), "test");
  final cm2 = CommandeFromJson(CommandeToJson(cm));
  print("${cm.date_emission}, ${cm2.date_emission}");

  final asString = jsonEncode(js);
  print(asString);

  LivraisonFromJson(jsonDecode(asString));
  print(jsonEncode(CommandeToJson(cm)));
}
