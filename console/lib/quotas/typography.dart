import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;
import './usage.dart';

class Typography extends StatelessWidget {
  static api.Quota zero = api.Quota();
  final api.Quota current;
  final Future<api.Quota> Function(api.Quota)? onChange;
  const Typography(this.current, {super.key, this.onChange});

  static FutureBuilder<api.Quota> future(
    Future<api.Quota> pending, {
    Future<api.Quota> Function(api.Quota)? onChange,
    api.Quota? defaults,
  }) {
    return ds.future(defaults ?? Typography.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        Typography(snapshot.data ?? Typography.zero, key: ValueKey(snapshot.data?.updatedAt), onChange: onChange),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final typography =
        "${ds.bytesx(current.consumed.toInt()).toIEC600272Format()} used of ${ds.bytesx(current.available().toInt()).toIEC600272Format()}";

    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing / 2,
      children: [Text(current.description), Usage(current), Text(typography)],
    );
  }
}
