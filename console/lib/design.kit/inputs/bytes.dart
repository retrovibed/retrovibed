import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:fixnum/fixnum.dart';
import '../bytesx.dart';
import '../theme.defaults.dart';

class Bytes extends StatefulWidget {
  final ValueChanged<Int64>? onChange;
  final Int64 value;
  final int magnitude;
  final InputDecoration? decoration;

  Bytes({
    super.key,
    this.onChange,
    int value = 0,
    int magnitude = bytesx.GiB,
    this.decoration,
  }) : value = Int64((value / magnitude).toInt()),
       magnitude = magnitude;

  @override
  State<Bytes> createState() => _ByteInputWidgetState(
    bytes: value,
    magnitude: magnitude,
    decoration: decoration,
  );
}

class _ByteInputWidgetState extends State<Bytes> {
  final TextEditingController _controller;
  int _magnitude;
  Int64 _bytes;
  final InputDecoration? _decoration;

  _ByteInputWidgetState({
    required Int64 bytes,
    required int magnitude,
    required InputDecoration? decoration,
  }) : _bytes = bytes,
       _magnitude = magnitude,
       _decoration = decoration,
       _controller = TextEditingController(text: bytes.toString());

  @override
  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    _controller.addListener(_updateBytes);
  }

  @override
  void dispose() {
    _controller.removeListener(_updateBytes);
    _controller.dispose();
    super.dispose();
  }

  void _updateBytes() {
    final text = _controller.text;
    final int? parsed = int.tryParse(text);

    if (parsed == null) {
      return;
    }

    final _parsed = Int64(parsed);
    if (_parsed == _bytes) {
      return;
    }

    setState(() {
      _bytes = _parsed;
    });

    final v = _bytes * _magnitude;
    widget.onChange?.call(v);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = Defaults.of(context);

    return Row(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing,
      children: [
        Flexible(
          child: TextField(
            controller: _controller,
            keyboardType: TextInputType.number, // Shows numeric keyboard
            inputFormatters: <TextInputFormatter>[
              FilteringTextInputFormatter.digitsOnly, // Allows only digits
            ],
            decoration: _decoration,
          ),
        ),
        DropdownButton<int>(
          value: _magnitude,
          underline: const SizedBox(),
          onChanged: (v) {
            setState(() {
              _magnitude = v ?? _magnitude;
            });
            final Int64 y = (_bytes * Int64(_magnitude));
            widget.onChange?.call(y);
          },
          items:
              [bytesx.KiB, bytesx.MiB, bytesx.GiB, bytesx.TiB].map((
                int magnitude,
              ) {
                return DropdownMenuItem<int>(
                  value: magnitude,
                  child: Text(bytesx.getName(magnitude)),
                );
              }).toList(),
        ),
      ],
    );
  }
}
