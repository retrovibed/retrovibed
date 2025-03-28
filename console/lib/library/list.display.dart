import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:desktop_drop/desktop_drop.dart';
import 'package:console/designkit.dart' as ds;
import 'package:console/media.dart' as media;
import 'package:console/mimex.dart' as mimex;
import './search.row.dart';

class AvailableListDisplay extends StatefulWidget {
  final media.FnMediaSearch search;
  final media.FnUploadRequest upload;
  final TextEditingController? controller;
  final FocusNode? focus;
  const AvailableListDisplay({
    super.key,
    this.search = media.media.get,
    this.upload = media.media.upload,
    this.controller,
    this.focus,
  });

  @override
  State<StatefulWidget> createState() => _AvailableListDisplay();
}

class _AvailableListDisplay extends State<AvailableListDisplay> {
  bool _loading = true;
  Widget? _player = null;
  ds.Error? _cause = null;
  media.MediaSearchResponse _res = media.media.response(
    next: media.media.request(limit: 32),
  );

  void refresh(media.MediaSearchRequest req) {
    widget
        .search(req)
        .then((v) {
          setState(() {
            _res = v;
            _loading = false;
          });

          widget.focus?.requestFocus();
          WidgetsBinding.instance.addPostFrameCallback((_) {
            final controller = widget.controller;
            if (controller == null) return;
            controller.selection = TextSelection.fromPosition(
              TextPosition(offset: controller.text.length),
            );
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
    refresh(_res.next);
  }

  @override
  Widget build(BuildContext context) {
    final upload = (DropDoneDetails v) {
      setState(() {
        _loading = true;
      });

      final multiparts = v.files.map((c) {
        return media.media.uploadable(c.path, c.name, c.mimeType!);
      });

      return Future.microtask(() {
        return Future.wait(
              multiparts.map((fv) {
                return fv.then((v) {
                  return widget
                      .upload((req) {
                        req..files.add(v);
                        return req;
                      })
                      .then((uploaded) {
                        setState(() {
                          _res.items.add(uploaded.media);
                        });
                      })
                      .catchError((cause) {
                        setState(() {
                          _cause = ds.Error.unknown(cause);
                        });
                      });
                });
              }),
            )
            .then((v) => ds.NullWidget)
            .catchError(ds.Error.unknown)
            .whenComplete(
              () => setState(() {
                _loading = false;
              }),
            );
      });
    };

    return ds.Overlay(
      overlay: _player,
      child: ds.Table(
        loading: _loading,
        cause: _cause,
        leading: SearchTray(
          controller: widget.controller,
          focus: widget.focus,
          onSubmitted: (v) {
            setState(() {
              _res.next.query = v;
              _res.next.offset = fixnum.Int64(0);
            });
            refresh(_res.next);
          },
          next: (i) {
            setState(() {
              _res.next.offset = i;
            });
            refresh(_res.next);
          },
          current: _res.next.offset,
          empty: fixnum.Int64(_res.items.length) < _res.next.limit,
          trailing: ds.FileDropWell(
            upload,
            child: IconButton(
              onPressed: () {},
              icon: Icon(Icons.file_upload_outlined),
            ),
          ),
          autofocus: true,
        ),
        children: _res.items,
        flex: 1,
        ds.Table.expanded<media.Media>(
          (v) => media.RowDisplay(
            media: v,
            leading: [Icon(mimex.icon(v.mimetype))],
            trailing: [media.ButtonShare(current: v)],
            onTap: media.PlayAction(context, v),
          ),
        ),
        empty: ds.FileDropWell(upload),
      ),
    );
  }
}
