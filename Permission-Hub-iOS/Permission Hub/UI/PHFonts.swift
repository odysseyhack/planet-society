//
//  PHFonts.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum PHFonts {

    static func bold(ofSize size: CGFloat = 13) -> UIFont {
        return UIFont(name: "Helvetica-Bold", size: size)!
    }

    static func regular(ofSize size: CGFloat = 13) -> UIFont {
        return UIFont(name: "Helvetica", size: size)!
    }
}
