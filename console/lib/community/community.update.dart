import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'community.pb.dart';
import 'community.edit.dart';

class CommunityUpdate extends StatefulWidget {
  final Community community;
  final Future<CommunityUpdateResponse> Function(Community) update;
  final void Function(Community) onUpdate;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Clip clipBehavior;
  final Duration debounce;

  const CommunityUpdate({
    super.key,
    required this.community,
    required this.update,
    required this.onUpdate,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.clipBehavior = Clip.none,
    this.debounce = const Duration(milliseconds: 500),
  });

  @override
  _CommunityUpdateState createState() => _CommunityUpdateState();
}

class _CommunityUpdateState extends State<CommunityUpdate> with ds.LoadingState {
  Community _community = Community();
  Timer _debounce = Timer(Duration.zero, () {});

  @override
  void initState() {
    super.initState();
    _debounce.cancel();
    _community = widget.community.deepCopy();
  }

  @override
  void dispose() {
    _debounce.cancel();
    super.dispose();
  }

  // _schedule coalesces rapid updates into a single _save.
  void _schedule() {
    _debounce.cancel();
    _debounce = Timer(widget.debounce, _save);
  }

  Future<void> _save() {
    _debounce.cancel();
    setState(() {
      loading = true;
      cause = ds.Error.zero;
    });

    return widget
        .update(_community)
        .then((response) => widget.onUpdate(response.community))
        .catchError((err) {
          setState(() {
            loading = false;
            cause = ds.Errors.httpauto(err, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((err) {
          setState(() {
            loading = false;
            cause = ds.Error.unknown(err, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: widget.padding ?? defaults.padding,
      margin: widget.margin,
      decoration: widget.decoration ?? const BoxDecoration(),
      constraints: widget.constraints,
      alignment: widget.alignment,
      clipBehavior: widget.clipBehavior,
      Column(
        spacing: defaults.spacing,
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          ds.Loading(
            cause: cause,
            CommunityEdit(
              community: _community,
              onChange: (c) {
                setState(() => _community = c);
                _schedule();
              },
              readOnly: true,
              autofocus: true,
            ),
          ),
        ],
      ),
    );
  }
}
