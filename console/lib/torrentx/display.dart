import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/torrentx/api.dart' as api;

export 'package:retrovibed/torrentx/api.dart';

class TorrentDisplay extends StatelessWidget {
  final api.TorrentInfoResponse current;
  final Color? background;

  const TorrentDisplay(this.current, {super.key, this.background});

  static Widget fromID(
    String id, {
    Key? key,
    Color? background,
    api.FnTorrentInfo get = api.info.get,
  }) {
    return Builder(
      key: key,
      builder: (context) {
        return FutureBuilder<api.TorrentInfoResponse>(
          initialData: api.TorrentInfoResponse(),
          future: get(id, options: [authn.request(authn.AuthzCache.meta(context))]),
          builder: (BuildContext ctx, AsyncSnapshot<api.TorrentInfoResponse> snapshot) {
            return ds.Loading(
              loading: !(snapshot.hasData || snapshot.hasError),
              cause: ds.Error.maybeErr(snapshot.error),
              snapshot.data == null ? SizedBox() : TorrentDisplay(snapshot.data!, background: background),
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

    return ds.Container(
      decoration: BoxDecoration(color: background),
      padding: defaults.padding,
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
          forms.Field(label: Text("source"), input: Text(details.source)),
          Transform.translate(
            child: forms.Checkbox(Text("private"), value: details.private),
            offset: Offset(-10.0, 0.0),
          ),
          forms.Field(label: Text("comment"), input: Text(meta.comment)),
          forms.Field(label: Text("encoding"), input: Text(meta.encoding)),
          forms.Field(label: Text("created by"), input: Text(meta.createdBy)),
          forms.Field(
            label: Text("trackers"),
            input: Text(meta.announceList.isEmpty ? '—' : meta.announceList.join(', ')),
          ),
          if (meta.urlList.isNotEmpty)
            forms.Field(label: Text("web seeds"), input: Text(meta.urlList.join(', '))),
          Text("Files", style: theme.textTheme.titleSmall),
          if (current.files.isEmpty)
            Text('—')
          else
            SizedBox(
              height: 220,
              child: ListView.builder(
                shrinkWrap: true,
                itemCount: current.files.length,
                itemBuilder: (context, i) {
                  final f = current.files[i];
                  return forms.Field(
                    input: Row(
                      children: [
                        Icon(Icons.insert_drive_file, size: 16),
                        SizedBox(width: 8),
                        Expanded(child: Text(f.path, overflow: TextOverflow.ellipsis, maxLines: 1)),
                        ds.Bytes(f.length),
                      ],
                    ),
                  );
                },
              ),
            ),
        ],
      ),
    );
  }
}
