import 'package:flutter/material.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;

class Usage extends StatelessWidget {
  static api.Quota zero = api.Quota();
  final api.Quota current;
  final double height;
  Usage(this.current, {super.key, this.height = 8.0});

  static FutureBuilder<api.Quota> future(
    Future<api.Quota> pending, {
    Future<api.Quota> Function(api.Quota)? onChange,
  }) {
    return ds.future(
      Usage.zero,
      pending.catchError((cause) {
        return api.Quota();
      }, test: httpx.ErrorsTest.err404),
      (snapshot) {
        return ds.ErrorScreen(
          Usage(snapshot.data ?? Usage.zero, key: ValueKey(snapshot.data?.updatedAt)),
          cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final available = current.available();
    return ds.Gauge(
      current.consumed.toInt() / available.toInt(),
      height: height,
    );
  }
}
