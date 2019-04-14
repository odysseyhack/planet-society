//
//  TransactionTermsViewController.swift
//  Permission Hub
//
//  Created by Corné on 14/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionTermsViewController: UIViewController {

    // MARK: - Private properties

    private let textView: UITextView = {

        let textView = UITextView()
        textView.font = PHFonts.regular()
        textView.textColor = PHColors.greyishBrown

        return textView
    }()

    let fileName: String

    // MARK: - Initialization

    init(fileName: String) {
        self.fileName = fileName

        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func loadView() {
        super.loadView()

        view = textView
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        let url = Bundle.main.url(forResource: fileName, withExtension: "txt")!
        let attrString = try! NSAttributedString(fileURL: url, options: [:], documentAttributes: nil)
        textView.attributedText = attrString
    }
}