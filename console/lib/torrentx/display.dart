import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/torrentx/api.dart' as api;

export 'package:retrovibed/torrentx/api.dart';

class TorrentDisplay extends StatelessWidget {
  final api.TorrentInfoResponse current;

  const TorrentDisplay(this.current, {super.key});

  static Widget fromID(
    String id, {
    Key? key,
    api.FnTorrentInfo get = api.info.get,
  }) {
    return Builder(
      key: key,
      builder: (context) {
        return FutureBuilder<api.TorrentInfoResponse>(
          initialData: api.TorrentInfoResponse(),
          // no torrent file has been fetched to local storage yet (e.g. an
          // undownloaded discovered item) - treat that as an empty result,
          // not an error, and let the display render disabled.
          future:
              get(
                id,
                options: [authn.request(authn.AuthzCache.meta(context))],
              ).catchError((e) {
                return api.TorrentInfoResponse();
              }, test: httpx.ErrorsTest.err404),
          builder: (BuildContext ctx, AsyncSnapshot<api.TorrentInfoResponse> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              TorrentDisplay(snapshot.data ?? api.TorrentInfoResponse()),
            );
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final meta = current.meta;
    final details = current.details;
    final ratio = details.downloaded == Int64.ZERO
        ? "-"
        : (details.uploaded.toInt() / details.downloaded.toInt()).toStringAsFixed(2);

    return Opacity(
      opacity: current.hasDetails() ? 1.0 : 0.4,
      child: ds.Container(
        padding: defaults.padding,
        margin: defaults.margin.copyWith(top: 0, bottom: 0),
        Column(
          mainAxisSize: MainAxisSize.min,
          mainAxisAlignment: MainAxisAlignment.start,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text("Torrent", style: theme.textTheme.titleSmall),
            forms.Field(
              label: Text("name"),
              input: Text(details.name, overflow: TextOverflow.ellipsis, maxLines: 1),
            ),
            forms.Field(label: Text("length"), input: ds.Bytes(details.length)),
            forms.Field(label: Text("downloaded"), input: ds.Bytes(details.downloaded)),
            forms.Field(label: Text("uploaded"), input: ds.Bytes(details.uploaded)),
            forms.Field(label: Text("ratio"), input: Text(ratio)),
            forms.Field(label: Text("source"), input: Text(details.source)),
            Transform.translate(
              child: forms.Checkbox(Text("private"), value: details.private),
              offset: Offset(-10.0, 0.0),
            ),
            forms.Field(label: Text("comment"), input: Text(meta.comment)),
            forms.Field(label: Text("encoding"), input: Text(meta.encoding)),
            forms.Field(label: Text("created by"), input: Text(meta.createdBy)),
            forms.Field(
              label: Text("files"),
              input: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: current.files
                    .map(
                      (f) => Row(
                        spacing: defaults.spacing,
                        children: [
                          ds.Bytes(f.length),
                          Expanded(child: Text(f.path, overflow: TextOverflow.ellipsis, maxLines: 1)),
                        ],
                      ),
                    )
                    .toList(),
              ),
            ),
            forms.Field(
              label: Text("trackers"),
              input: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: meta.announceList.map((v) => Text(v, overflow: TextOverflow.ellipsis, maxLines: 1)).toList(),
              ),
            ),
            forms.Field(
              label: Text("web seeds"),
              input: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: meta.urlList.map((v) => Text(v, overflow: TextOverflow.ellipsis, maxLines: 1)).toList(),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
