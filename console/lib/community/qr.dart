import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:retrovibed/community/community.pb.dart';

String encodeQRPayload(Community community, {String attribution = ''}) {
  return jsonEncode(<String, dynamic>{
    'community': community.toProto3Json(),
    'attribution': attribution,
  });
}

(Community?, String) decodeQRPayload(String data) {
  try {
    final payload = jsonDecode(data);
    if (payload is! Map<String, dynamic>) return (null, '');
    if (payload.containsKey('community')) {
      final community = Community.create()
        ..mergeFromProto3Json(payload['community']);
      final attribution = payload['attribution'] as String? ?? '';
      return (community, attribution);
    }
    final community = Community.create()..mergeFromProto3Json(payload);
    return (community, '');
  } catch (e, s) {
    debugPrint('decodeQRPayload failed: $e\n$s');
    return (null, '');
  }
}
