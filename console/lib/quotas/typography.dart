import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './api.dart' as api;
import './usage.dart';

class Typography extends StatelessWidget {
  static api.Quota zero = api.Quota();
  final api.Quota current;
  final String label;
  final Future<api.Quota> Function(api.Quota)? onChange;
  const Typography(this.current, {super.key, this.label = "", this.onChange});

  static FutureBuilder<api.Quota> future(
    Future<api.Quota> pending, {
    Future<api.Quota> Function(api.Quota)? onChange,
    api.Quota? defaults,
    String label = "",
  }) {
    return ds.future(defaults ?? Typography.zero, pending, (snapshot) {
      return ds.ErrorScreen(
        Typography(snapshot.data ?? Typography.zero, key: ValueKey(snapshot.data?.updatedAt), label: label, onChange: onChange),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final title = label.isEmpty ? current.description : label;
    final typography =
        "${ds.bytesx(current.consumed.toInt()).toIEC600272Format()} used of ${ds.bytesx(current.available().toInt()).toIEC600272Format()}";

    return Column(
      mainAxisSize: MainAxisSize.min,
      mainAxisAlignment: MainAxisAlignment.start,
      crossAxisAlignment: CrossAxisAlignment.start,
      spacing: defaults.spacing / 2,
      children: [Text(title), Usage(current), Text(typography)],
    );
  }
}
