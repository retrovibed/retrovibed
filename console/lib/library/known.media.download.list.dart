import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import './known.media.locator.dart';
import './known.media.source.dart';

class KnownMediaDownloadList extends StatelessWidget {
  final Widget leading;
  final List<api.Known> children;
  final Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate;

  const KnownMediaDownloadList({
    super.key,
    this.children = const [],
    this.leading = const SizedBox(),
    this.locate = api.locate.create,
  });

  static Widget query(
    Future<List<api.Known>> Function() query, {
    Key? key,
    Future<api.LocateCreateResponse> Function(api.Locate req, {List<httpx.Option> options}) locate = api.locate.create,
    Widget leading = const SizedBox(),
  }) {
    return FutureBuilder<List<api.Known>>(
      initialData: [],
      future: query(),
      builder: (BuildContext ctx, AsyncSnapshot<List<api.Known>> snapshot) {
        return ds.Loading(
          loading: !(snapshot.hasData || snapshot.hasError),
          cause: ds.Error.maybeErr(snapshot.error),
          KnownMediaDownloadList(
            leading: leading,
            children: snapshot.data ?? [],
            locate: locate,
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return ds.Grid(
      leading: [leading],
      children: children,
      (context, v) => KnownMediaLocator(
        v,
        locate: locate,
        trailing: Center(child: KnownMediaSource(v, height: 24)),
      ),
    );
  }
}
