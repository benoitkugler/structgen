import 'dart:convert';

import 'gen.dart';

main(List<String> args) {
  final m =
      Model(Concret1([1, 2, -1], 4), 8, [Concret2(0.4), Concret2(0.8)], {3: 4});
  final json = modelToJson(m);
  final s = jsonEncode(json);
  print(s);

  final decoded = jsonDecode(s);
  final s2 = jsonEncode(modelToJson(modelFromJson(decoded)));

  if (s != s2) {
    throw ("inconstistent roundtrip");
  }
  print("OK");
}
