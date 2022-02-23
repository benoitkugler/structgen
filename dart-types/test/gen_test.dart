import 'dart:convert';

import 'gen.dart';

main(List<String> args) {
  final m = model(concret1([1, 2, -1], 4), 8, [concret2(0.4), concret2(0.8)]);
  final json = modelToJson(m);
  final s = jsonEncode(json);
  print(s);

  final decoded = jsonDecode(s);
  final s2 = jsonEncode(modelToJson(modelFromJson(decoded)));

  if (s != s2) {
    throw ("inconstistent roundtrip");
  }
}
