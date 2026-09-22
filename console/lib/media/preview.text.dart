import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import './api.dart' as api;
import './media.pb.dart';
import './preview.unsupported.dart';

// reads the head of the file rather than the whole thing. the daemon serves Range because
// blockcache.File seeks, so this pulls one block from cloud storage instead of hydrating
// the entire object to show the first screenful.
class PreviewText extends StatefulWidget {
  static const int limit = 64 * 1024;

  final Media current;
  final api.FnMediaDownload download;

  const PreviewText({
    super.key,
    required this.current,
    this.download = api.media.download,
  });

  @override
  State<StatefulWidget> createState() => _PreviewText();
}

class _PreviewText extends State<PreviewText> with ds.LoadingState {
  bool _truncated = false;
  String _content = "";

  @override
  void initState() {
    super.initState();
    ds.postframe(read);
  }

  Future<void> read() {
    return httpx
        .withRetry(
          () => widget.download(
            widget.current.id,
            options: [
              httpx.Request.header("Range", "bytes=0-${PreviewText.limit - 1}"),
              authn.request(authn.AuthzCache.meta(context)),
            ],
          ),
        )
        .then((v) => v.stream.bytesToString())
        .then((body) {
          setState(() {
            _content = body;
            // a full read means the file is at least this long, so the reader is looking
            // at a head rather than the document.
            _truncated = body.length >= PreviewText.limit;
            loading = false;
          });
        })
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unauthorized(cause, onTap: reseterr);
            loading = false;
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
            loading = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final theme = Theme.of(context);

    return ds.Loading(
      cause: cause,
      loading: loading,
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Padding(
            padding: defaults.padding,
            child: SelectableText(
              _content,
              style: theme.textTheme.bodySmall?.copyWith(fontFamily: "monospace"),
            ),
          ),
          Visibility(
            visible: _truncated,
            child: PreviewUnsupported(
              current: widget.current,
              description: const Text("preview truncated, download to read the rest"),
            ),
          ),
        ],
      ),
    );
  }
}
