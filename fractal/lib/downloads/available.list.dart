import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as ds;
import 'package:fractal/media.dart' as media;

class AvailableListDisplay extends StatefulWidget {
  final media.FnMediaSearch search;
  const AvailableListDisplay({
    super.key,
    this.search = media.discovered.available,
  });

  @override
  State<StatefulWidget> createState() => _AvailableListDisplay();
}

class _AvailableListDisplay extends State<AvailableListDisplay> {
  bool _loading = true;
  ds.Error? _cause = null;
  media.MediaSearchResponse _res = media.media.response(
    next: media.media.request(limit: 32),
  );

  void refresh() {
    widget
        .search(_res.next)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });
        })
        .catchError((e) {
          setState(() {
            _cause = ds.Error.unknown(e);
            _loading = false;
          });
        });
  }

  @override
  void initState() {
    super.initState();
    refresh();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table(
      loading: _loading,
      cause: _cause,
      children: _res.items,
      leading: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Expanded(
                child: TextField(
                  decoration: InputDecoration(
                    hintText: "search available content",
                  ),
                  onChanged:
                      (v) => setState(() {
                        _res.next.query = v;
                      }),
                  onSubmitted: (v) {
                    setState(() {
                      _res.next.offset = fixnum.Int64(0);
                    });
                    refresh();
                  },
                ),
              ),
              IconButton(
                onPressed:
                    _res.next.offset > 0
                        ? () {
                          setState(() {
                            _res.next.offset -= 1;
                          });
                          refresh();
                        }
                        : null,
                icon: Icon(Icons.arrow_left),
              ),
              IconButton(
                onPressed:
                    _res.items.isNotEmpty
                        ? () {
                          setState(() {
                            _res.next.offset += 1;
                          });
                          refresh();
                        }
                        : null,
                icon: Icon(Icons.arrow_right),
              ),
            ],
          ),
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [Spacer(), Text("description"), Spacer()],
          ),
        ],
      ),
      (v) => media.RowDisplay(
        media: v,
        onTap:
            () => media.discovered
                .download(v.id)
                .then((v) {
                  refresh();
                })
                .catchError((cause) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text("Failed to download: $cause")),
                  );
                  return null;
                }),
      ),
    );
  }
}
