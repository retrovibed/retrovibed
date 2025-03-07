import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;

class Table<T> extends StatelessWidget {
  final Widget Function(T i) render;
  final List<T> children;
  final Widget empty;
  final Widget leading;
  final Widget trailing;
  final bool loading;
  final ds.Error? cause;

  const Table(
    this.render, {
    super.key,
    this.leading = const SizedBox(),
    this.trailing = const SizedBox(),
    this.empty = const SizedBox(),
    this.children = const [],
    this.loading = false,
    this.cause = null,
  });

  @override
  Widget build(BuildContext context) {
    final content =
        this.children.length == 0
            ? this.empty
            : ListView(
              shrinkWrap: true,
              children: this.children.map(this.render).toList(),
            );
    return ds.Loading(
      loading: loading,
      cause: cause,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        mainAxisSize: MainAxisSize.max,
        children: [leading, Expanded(child: content), trailing],
      ),
    );
  }
}
