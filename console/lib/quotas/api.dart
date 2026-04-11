import 'dart:convert';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './meta.quota.pb.dart';

export './meta.quota.pb.dart';

abstract class api {
  static Future<QuotaSearchResponse> search({
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.metaendpoint(), "/q/", {}), options: options).then((v) {
      return Future.value(
        QuotaSearchResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
      );
    });
  }

  static Future<QuotaFindResponse> get(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .get(
          Uri.https(httpx.metaendpoint(), "/q/${id}", {}),
          options: options,
        )
        .then((v) {
          return Future.value(
            QuotaFindResponse.create()..mergeFromProto3Json(jsonDecode(v.body)),
          );
        });
  }
}

extension Exts on Quota {
  ds.Int64 remaining() {
    return (this.credits + this.maximum + this.granted - this.consumed - this.reserved).clamp(
      ds.Int64.ZERO,
      ds.Int64.MAX_VALUE,
    );
  }

  ds.Int64 available() {
    return (this.credits + this.maximum + this.granted).clamp(ds.Int64.ZERO, ds.Int64.MAX_VALUE);
  }
}
