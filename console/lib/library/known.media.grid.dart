import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'known.media.card.dart';
import './api.dart' as api;

class KnownMediaGrid extends StatelessWidget {
  final List<api.Known> children;
  const KnownMediaGrid({
    super.key,
    this.children = const [],
  });

  @override
  Widget build(BuildContext context) {
    return ds.Grid(
      children: children,
      (context, current) => KnownMediaCard(current),
    );
  }
}
