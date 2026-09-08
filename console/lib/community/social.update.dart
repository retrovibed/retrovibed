import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'community.social.pb.dart';
import 'social.edit.dart';

class SocialUpdate extends StatefulWidget {
  final CommunityPublisher compub;
  final Future<CommunityPublisherUpdateResponse> Function(CommunityPublisher) update;
  final void Function(CommunityPublisher) onUpdate;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Clip clipBehavior;
  final List<Widget> actions;
  final Duration debounce;

  const SocialUpdate({
    super.key,
    required this.compub,
    required this.update,
    required this.onUpdate,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.clipBehavior = Clip.none,
    this.actions = const [],
    this.debounce = const Duration(milliseconds: 500),
  });

  @override
  _SocialUpdateState createState() => _SocialUpdateState();
}

class _SocialUpdateState extends State<SocialUpdate> with ds.LoadingState {
  CommunityPublisher _compub = CommunityPublisher();
  Timer _debounce = Timer(Duration.zero, () {});

  @override
  void initState() {
    super.initState();
    _debounce.cancel();
    _compub = widget.compub.deepCopy();
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
    setState(() {
      loading = true;
      cause = ds.Error.zero;
    });

    return widget
        .update(_compub)
        .then((response) => widget.onUpdate(response.compub))
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
            SocialEdit(
              compub: _compub,
              onChange: (c) {
                setState(() => _compub = c);
                _schedule();
              },
              autofocus: true,
              actions: widget.actions,
            ),
          ),
        ],
      ),
    );
  }
}
