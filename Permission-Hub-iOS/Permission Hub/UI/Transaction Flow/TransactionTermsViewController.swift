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

    private let notificationCell: TransactionNotificationTableViewCell = {

        let cell = TransactionNotificationTableViewCell()
        cell.translatesAutoresizingMaskIntoConstraints = false
        cell.configure(
            withType: .warning,
            andText: "You are allowed to cancel this agreement within 14 days.")

        return cell
    }()

    private let textView: UITextView = {

        let textView = UITextView()
        textView.translatesAutoresizingMaskIntoConstraints = false
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

    override func viewDidLoad() {
        super.viewDidLoad()

        view.addSubview(notificationCell)
        view.addSubview(textView)

        notificationCell.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        notificationCell.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        notificationCell.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        notificationCell.heightAnchor.constraint(equalToConstant: 44).isActive = true

        textView.topAnchor.constraint(equalTo: notificationCell.bottomAnchor).isActive = true
        textView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        textView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        textView.bottomAnchor.constraint(equalTo: view.bottomAnchor).isActive = true

        let url = Bundle.main.url(forResource: fileName, withExtension: "txt")!
        let attrString = try! NSAttributedString(fileURL: url, options: [:], documentAttributes: nil)
        textView.attributedText = attrString
    }
}
