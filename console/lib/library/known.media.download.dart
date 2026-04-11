import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'known.media.card.dart';
import './api.dart' as api;
import './known.media.source.dart';

class KnownMediaDownload extends StatefulWidget {
  final Widget leading;
  final List<api.Known> children;
  final Future<api.Known> Function(api.Known k) onTap;
  final Future<api.Known> Function(api.Known k)? onDoubleTap;
  const KnownMediaDownload({
    super.key,
    this.children = const [],
    this.leading = const SizedBox(),
    required this.onTap,
    this.onDoubleTap,
  });

  static Widget query(
    Future<List<api.Known>> Function() query, {
    Key? key,
    required Future<api.Known> Function(api.Known k) onTap,
    Future<api.Known> Function(api.Known k)? onDoubleTap,
    Widget leading = const SizedBox(),
  }) {
    return FutureBuilder<List<api.Known>>(
      initialData: [],
      future: query(),
      builder: (BuildContext ctx, AsyncSnapshot<List<api.Known>> snapshot) {
        return ds.Loading(
          loading: !(snapshot.hasData || snapshot.hasError),
          cause: ds.Error.maybeErr(snapshot.error),
          KnownMediaDownload(
            leading: leading,
            children: snapshot.data ?? [],
            onTap: onTap,
            onDoubleTap: onDoubleTap,
          ),
        );
      },
    );
  }

  @override
  State<StatefulWidget> createState() => _KnownMediaDownload();
}

class _KnownMediaDownload extends State<KnownMediaDownload> {
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Grid(
      leading: [widget.leading],
      children: widget.children,
      (context, v) => KnownMediaCard(
        v,
        icon: Icons.search,
        onTap: () => widget.onTap(v),
        onDoubleTap: widget.onDoubleTap == null ? null : () => widget.onDoubleTap!(v),
        trailing: Center(child: KnownMediaSource(v, height: 24)),
      ),
    );
  }
}
