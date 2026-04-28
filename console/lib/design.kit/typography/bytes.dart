import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/design.kit/bytesx.dart';

class Bytes extends StatelessWidget {
  static String FormatIEC600272(Int64 v) => bytesx(v.toInt()).toIEC600272Format();
  static String FormatSI(Int64 v) => bytesx(v.toInt()).toSIFormat();
  final Int64 duration;
  final String Function(Int64 v) format;

  const Bytes(this.duration, {super.key, this.format = Bytes.FormatIEC600272});

  factory Bytes.SI(Int64 v) {
    return Bytes(v, format: Bytes.FormatSI);
  }

  factory Bytes.Int(int v) {
    return Bytes(Int64(v));
  }

  @override
  Widget build(BuildContext context) {
    return Text(format(duration));
  }
}
