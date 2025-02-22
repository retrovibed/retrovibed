import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;

class Table extends StatelessWidget {
  final List<Widget> children;
  final Widget leading;
  final Widget trailing;
  final bool loading;
  final ds.Error? cause;

  const Table({
    super.key,
    this.leading = const SizedBox(),
    this.trailing = const SizedBox(),
    this.children = const [],
    this.loading = false,
    this.cause = null,
  });

  @override
  Widget build(BuildContext context) {
    return ds.Loading(
      loading: loading,
      cause: cause,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          leading,
          Expanded(child: ListView(shrinkWrap: true, children: children)),
          trailing,
        ],
      ),
    );
  }
}
